package main

import (
	"os"
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// init configures the zerolog logger before the main function runs.
func init() {
	// --- Zerolog Configuration ---
	// Set the global log level. For a migration script
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Configure Zerolog's global logger to output JSON to os.Stdout.
	// and enable console logging for better readability
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		With().Timestamp().Caller().Logger()
}

// main is the entry point for the database migration script.
// It connects to the PostgreSQL database using the DATABASE_URL environment variable,
// then executes all pending schema migrations defined within the `pkg/database` package.
// This script is crucial for setting up and updating the database schema.
// A fatal error occurs if DATABASE_URL is not provided or if the database connection fails.
func main() {
	// --- Database Migration Process ---
	log.Info().Str("event", "migration_start").Msg("Starting database migration process")

	// Check if DATABASE_URL is set
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal().Str("event", "env_var_missing").Msg("Error: DATABASE_URL environment variable is not set. Cannot run migrations.")
	}

	log.Info().Str("event", "env_var_checked").Msg("DATABASE_URL environment variable found.")

	// Connect to the database
	log.Info().Str("event", "db_connect_attempt").Msg("Attempting to connect to database for migration...")

	db, err := database.ConnectMigrateDB(dbURL) // Using a dedicated function for migration DB connection for direct connectivity
	if err != nil {
		// Use Fatal for critical errors during database connection
		log.Fatal().Err(err).Str("event", "db_connect_failure").Msg("Error connecting to database for migration")
	}
	log.Info().Str("event", "db_connected").Msg("Successfully connected to database for migration.")

	// Defer closing the database connection
	defer func() {
		sqlDB, closeErr := db.DB()
		if closeErr == nil {
			if err := sqlDB.Close(); err != nil {
				// Use Error for non-fatal errors that should still be logged
				log.Error().Err(err).Str("event", "db_close_failure").Msg("Error closing database connection")
			} else {
				log.Info().Str("event", "db_closed").Msg("Database connection successfully closed.")
			}
		} else {
			log.Error().Err(closeErr).Str("event", "get_sql_db_failure").Msg("Error getting underlying SQL DB for closing")
		}
	}()

	// Run the migrations
	log.Info().Bool("isFullMigration", true).Str("event", "migrations_run_start").Msg("Running database migrations...")
	// if isFullMigration is true, this function will first **DROP ALL TABLES** corresponding to the registered models before re-creating them.
	database.Migrate(db, true)

	log.Info().Str("event", "migration_complete").Msg("Database migration process completed successfully!")
}
