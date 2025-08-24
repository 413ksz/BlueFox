package databaseerrorhelper

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func GetDatabaseErrorMessage(result *gorm.DB) *models.CustomError {
	if result.Error != nil {
		// Attempt to unwrap the error to check for specific database driver errors
		if pgErr, ok := result.Error.(*pgconn.PgError); ok {
			// Check for PostgreSQL unique violation error code (23505)
			if pgErr.Code == "23505" {
				// Handle specific unique constraint violations
				if pgErr.ConstraintName == "uni_users_email" {
					customError := models.NewCustomError(models.ERROR_CODE_CONFLICT, "A user with similar email already exists", result.Error, nil)
					return customError
				}
				// Fallback for any other unique constraint violation not specifically handled
				customError := models.NewCustomError(models.ERROR_CODE_CONFLICT, "An unspecific unique constraint violation occurred", result.Error, nil)
				return customError
			}
		}

		// Fallback for any other database errors (e.g., connection issues, other integrity errors)
		customError := models.NewCustomError(models.ERROR_CODE_INTERNAL_SERVER, "An unspecific database error occurred", result.Error, nil)
		return customError
	}

	return nil
}
