package database

import (
	"fmt"
	"os"
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global GORM database connection pool for the primary application operations.
var DB *gorm.DB

// InitAppDB initializes and configures the global GORM database connection pool (DB)
// for the primary application operations.
//
// It retrieves the PostgreSQL connection string from the "DATABASE_URL" environment variable.
// This function must be called once at application startup (e.g., in main.go)
// to ensure a ready and optimized database connection is available globally.
//
// The connection pool is configured with:
//   - MaxIdleConns: 5 	(Maximum number of connections in the idle connection pool)
//   - MaxOpenConns: 20 	(Maximum number of open connections to the database)
//   - ConnMaxLifetime: 1 minute (Maximum amount of time a connection may be reused)
//
// Returns:
//
//	error: An error if the "DATABASE_URL" environment variable is not set,
//	 	or if the connection to the database fails.
func InitAppDB() error {
	// --- Database Configuration ---
	log.Info().Msg("Initializing global database connection...")

	// Retrieve the database URL from environment variables.
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Error().Msg("DATABASE_URL environment variable is not set")
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	var err error
	// Open a new GORM database connection using the PostgreSQL driver.
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the underlying sql.DB object to configure connection pooling.
	sqlDB, err := DB.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get underlying SQL DB for app")
		return fmt.Errorf("failed to get underlying SQL DB for app: %w", err)
	}

	// Configure the connection pool to manage database connections efficiently.
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(1 * time.Minute)

	log.Info().Msg("Global application database connection established and pooled.")
	return nil
}

// ConnectMigrateDB establishes a dedicated GORM database connection for migration purposes.
// This function is separate from InitAppDB to allow for distinct connection
// parameters or lifecycle for migration operations, particularly useful in CLI tools
// or one-off scripts.
//
// It connects to the PostgreSQL database using the provided `dbURL` and performs
// a `Ping` to ensure the database is reachable before returning the connection.
// A connection pool is also configured specifically for this migration connection.
//
// Parameters:
//
//	dbURL: The PostgreSQL connection string (e.g., "postgres://user:password@localhost:5432/dbname?sslmode=disable").
//
// Returns:
//
//	*gorm.DB: A pointer to the GORM database connection instance configured for migration.
//	error: An error if the connection fails, the underlying SQL DB cannot be retrieved,
//	 	or if the database is unreachable via a ping.
func ConnectMigrateDB(dbURL string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database for migration")
		return nil, fmt.Errorf("failed to connect to database for migration: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get underlying SQL DB for migration")
		return nil, fmt.Errorf("failed to get underlying SQL DB for migration: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Error().Err(err).Msg("Failed to ping database for migration")
		return nil, fmt.Errorf("failed to ping database for migration: %w", err)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	log.Info().Msg("Migration database connection established.")
	return db, nil
}

// Migrate performs GORM's auto-migrations on the provided database instance.
// This function maps your Go struct models to database tables and columns.
//
// How GORM AutoMigrate works:
//   - Creates tables for models that do not exist.
//   - Adds missing columns to existing tables.
//   - Adds missing indexes and foreign key constraints.
//
// Important Note on Destructive Migrations:
// GORM's AutoMigrate is non-destructive by default. It WILL NOT:
//   - Delete columns, tables, or indexes that are no longer defined in your models.
//   - Modify existing column types.
//
// For such destructive or complex schema changes, consider using a dedicated
// versioned migration tool (e.g., "golang-migrate/migrate") or manual SQL scripts.
//
// If `isFullMigration` is true, this function will first **DROP ALL TABLES**
// corresponding to the registered models before re-creating them. This is
// primarily intended for **development environments** to ensure a clean slate
// and should be used with extreme caution in production, as it results in
// **DATA LOSS**.
//
// Parameters:
//
//	db: The GORM database instance on which to run the migrations.
//	isFullMigration: A boolean flag. If true, all specified tables will be
//	 	dropped before migration. Use only for development.
func Migrate(db *gorm.DB, isFullMigration bool) {
	log.Info().Msg("Starting database auto-migration...")
	if isFullMigration {
		log.Warn().Msg("Performing full migration: Dropping all specified tables. DATA LOSS WILL OCCUR!")
		db.Migrator().DropTable(
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
		log.Info().Msg("All specified tables dropped for full migration.")
	}
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
		log.Fatal().Err(err).Msg("Failed to auto-migrate database")
	}
	log.Info().Msg("Database auto-migration completed successfully!")
}
