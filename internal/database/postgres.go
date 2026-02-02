package database

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

type Config struct {
	DSN string
}

func New(cfg Config) (*DB, error) {
	dsn := cfg.DSN

	// Append prefer_simple_protocol to DSN to avoid prepared-statement cache issues
	if u, err := url.Parse(dsn); err == nil {
		q := u.Query()
		if q.Get("prefer_simple_protocol") == "" {
			q.Set("prefer_simple_protocol", "true")
			u.RawQuery = q.Encode()
			dsn = u.String()
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Disable PrepareStmt to avoid prepared-statement caching conflicts
		// (helpful when using connection poolers like pgbouncer/supabase pooler)
		PrepareStmt: false,
		Logger:      logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Successfully connected to PostgreSQL (Supabase)")

	return &DB{db}, nil
}

func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (db *DB) HealthCheck() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
