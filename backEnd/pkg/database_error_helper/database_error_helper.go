package databaseerrorhelper

import (
	"database/sql"
	"errors"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/jackc/pgx/v5/pgconn"
)

func GetDatabaseErrorMessage(databaseError error) *models.CustomError {
	if databaseError == nil {
		return nil
	}
	if errors.Is(databaseError, sql.ErrNoRows) {
		return models.NewCustomError(models.ERROR_CODE_NOT_FOUND, "Resource not found within the database", databaseError, nil)
	}
	var pgErr *pgconn.PgError
	if errors.As(databaseError, &pgErr) {
		if pgErr.Code == "23505" {
			if pgErr.ConstraintName == "user_email_key" {
				return models.NewCustomError(models.ERROR_CODE_CONFLICT, "Email already exists", databaseError, nil)
			} else {
				return models.NewCustomError(models.ERROR_CODE_CONFLICT, "An unspecified unique constraint was violated", databaseError, nil)
			}
		}
		return models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Unknown database error occurred", databaseError, nil)
	}
	return models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "Unknown not database specific error occurred", databaseError, nil)
}
