package migrations

import (
	"gocats/internal/database"
	"gocats/internal/models"
	"log"
	"strings"
)

func RunMigrations(db *database.DB) error {
	migrator := database.NewMigrator(db.DB)

	// List of models to migrate
	models := []interface{}{
		&models.Category{},          // Ensure Category is migrated before Product
		&models.Product{},           // Has a foreign key to Category
		&models.Transaction{},       // Transaction table
		&models.TransactionDetail{}, // Has foreign keys to Transaction and Product
	}

	if err := migrator.AutoMigrate(models...); err != nil {
		// Don't exit the program for migration errors. Log and continue.
		// Ignore common idempotent errors like "relation already exists" or prepared statement conflicts.
		msg := err.Error()
		if strings.Contains(msg, "already exists") || strings.Contains(msg, "stmtcache") {
			log.Printf("Migration warning (ignored): %v", err)
			return nil
		}

		log.Printf("Migration error (non-fatal): %v", err)
		return nil
	}

	log.Println("All Migrations completed")
	return nil
}
