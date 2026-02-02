package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) AutoMigrate(models ...interface{}) error {
	log.Println("🔄 Running database migrations...")

	if err := m.db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✅ Migrations completed successfully")
	return nil
}

func (m *Migrator) Migrate(models ...interface{}) error {
	log.Println("Running Database migrations ...")
	if err := m.db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("Migration failed: %w", err)
	}
	log.Println("Database migrations completed successfully")
	return nil
}

func (m *Migrator) DropTables(models ...interface{}) error {
	log.Println("🗑️  Dropping tables...")

	if err := m.db.Migrator().DropTable(models...); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	log.Println("Tables dropped successfully")
	return nil
}
