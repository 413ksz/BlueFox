package main

import (
	"log"
	"os"

	"github.com/413ksz/BlueFox/backEnd/pkg/database"
)

// main is the entry point for the database migration script.
// It connects to the PostgreSQL database using the DATABASE_URL environment variable,
// then executes all pending schema migrations defined within the `pkg/database` package.
//
// This script is crucial for setting up and updating the database schema.
// A fatal error occurs if DATABASE_URL is not provided or if the database connection fails.
func main() {
	log.Println("Starting database migration process...")

	// Check if DATABASE_URL is set
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("Error: DATABASE_URL environment variable is not set. Cannot run migrations.")
	}

	// Connect to the database
	db, err := database.ConnectMigrateDB(dbURL) // Using a dedicated function for migration DB connection for direct connectivity
	if err != nil {
		log.Fatalf("Error connecting to database for migration: %v", err)
	}
	defer func() {
		sqlDB, closeErr := db.DB()
		if closeErr == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Error closing database connection: %v", err)
			}
		}
	}()

	// Run the migrations
	database.Migrate(db)

	log.Println("Database migration process completed successfully!")
}
