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

// InitAppDB initializes the global database connection for your application's API handlers.
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

func CloseAppDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Error closing global database connection: %v", err)
			} else {
				log.Println("Global database connection closed.")
			}
		}
	}
}

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
