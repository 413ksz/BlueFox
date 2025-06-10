package database

import (
	"fmt"
	"log"
	"os"

	"github.com/413ksz/BlueFox/backEnd/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes and returns a GORM DB instance
func InitDB() (*gorm.DB, error) {
	// Get database URL from environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Default to SQLite if DATABASE_URL is not set
		dbURL = ""
		log.Println("DATABASE_URL environment variable not set.")
	}

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}

// AutoMigrate migrates the database schema for the given models.
// This is convenient for development but for production, consider explicit migrations.
// myapp/database/database.go
func Migrate(db *gorm.DB) {
	log.Println("Starting database auto-migration...")
	err := db.AutoMigrate(
		&models.User{},
		&models.Message{},
		&models.UserFriendConnect{},
		&models.Server{},
		&models.Channel{},
		&models.ServerUserConnect{},
		&models.MediaAsset{},
		&models.MessageAttachment{},
		// Add any new top-level models here.
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Println("Database auto-migration completed successfully!")
}
