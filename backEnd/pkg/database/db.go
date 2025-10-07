package database

import (
	"context"
	"os"
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// DB encapsulates a sql.DB database connection pool.
type DB struct {
	DB *pgxpool.Pool
}

// NewDB creates a new DB instance with the provided *sql.DB database connection.
func NewDB(db *pgxpool.Pool) *DB {
	return &DB{
		DB: db,
	}
}

// NewConnectionPool creates a new sql.DB database connection and returns a pointer to it.
// It also returns an error if the connection fails.
func NewConnectionPool() (*pgxpool.Pool, *models.CustomError) {
	log.Info().
		Str("component", "database").
		Str("event", "app_db_init_start").
		Msg("Initializing global database connection for sqlc")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Error().
			Str("component", "database").
			Str("event", "app_db_env_var_missing").
			Msg("DATABASE_URL environment variable is not set")
		return nil, models.NewCustomError(models.ERROR_CODE_INITIALIZE_ERROR, "DATABASE_URL environment variable is not set", nil, nil)
	}
	config, condigErr := pgxpool.ParseConfig(dbURL)
	if condigErr != nil {
		return nil, models.NewCustomError(models.ERROR_CODE_INITIALIZE_ERROR, "failed to parse database config", condigErr, nil)
	}
	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = time.Minute * 1
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, poolErr := pgxpool.NewWithConfig(ctx, config)
	if poolErr != nil {
		log.Error().
			Err(poolErr).
			Str("component", "database").
			Str("event", "app_db_connect_failure").
			Msg("Failed to connect to database")
		return nil, models.NewCustomError(models.ERROR_CODE_INITIALIZE_ERROR, "failed to connect to database using pgxpool", poolErr, nil)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Error().
			Err(err).
			Str("component", "database").
			Str("event", "app_db_ping_failure").
			Msg("Failed to ping database")
		return nil, models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "failed to ping database", err, nil)
	}

	log.Info().
		Str("component", "database").
		Str("event", "app_db_init_success").
		Int("max_idle_conns", 5).
		Int("max_open_conns", 20).
		Str("conn_max_lifetime", "1m").
		Msg("Global application database connection established and pooled.")

	return pool, nil
}
