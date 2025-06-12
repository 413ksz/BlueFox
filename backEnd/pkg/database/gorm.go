package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitAppDB initializes the global GORM database connection pool for the main application.
// It retrieves the database connection string from the "DATABASE_URL" environment variable.
// This function should be called once at application startup (e.g., in an init() function)
// to ensure a ready and configured database connection is available globally.
//
// Returns:
//
//	error: An error if the DATABASE_URL is not set or if the connection fails.
func InitAppDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB for app: %w", err)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	log.Println("Global application database connection established and pooled.")
	return nil
}

// ConnectMigrateDB establishes a connection to the PostgreSQL database using the provided
// database URL and configures a connection pool. It ensures the database is reachable by
// pinging it and sets up a connection pool with specified parameters.
//
// Parameters:
//
//	dbURL: A string representing the connection string for the PostgreSQL database.
//
// Returns:
//
//	*gorm.DB: A pointer to the GORM database connection instance.
//	error: An error if the connection fails, the underlying SQL DB cannot be retrieved,
//	       or if the database is unreachable.
func ConnectMigrateDB(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database for migration: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL DB for migration: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database for migration: %w", err)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	log.Println("Migration database connection established")
	return db, nil
}

// Migrate performs the GORM auto-migrations on the provided database instance.
// This function iterates through all registered GORM models and automatically:
// - Creates new tables for models that do not yet exist in the database.
// - Adds missing columns to existing tables.
// - Adds missing indexes and foreign key constraints.
//
// Important Note: GORM's AutoMigrate is non-destructive. It will NOT
// delete columns, tables, or indexes that are no longer defined in your models.
// For complex schema changes or deletions, a dedicated versioned migration tool
// (e.g., golang-migrate/migrate) or manual SQL scripts are typically required.
//
// Parameters:
//
//	db: The GORM database instance on which to run the migrations.
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
