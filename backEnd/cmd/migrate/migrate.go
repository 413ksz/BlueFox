package main

import (
	"log"
	"os"

	"github.com/413ksz/BlueFox/backEnd/pkg/database"
)

func main() {
	log.Println("Starting database migration process...")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("Error: DATABASE_URL environment variable is not set. Cannot run migrations.")
	}

	// Connect to the database
	db, err := database.ConnectMigrateDB(dbURL) // Using a dedicated function for migration DB connection
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
