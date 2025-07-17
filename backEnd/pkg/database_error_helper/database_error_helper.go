package databaseerrorhelper

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/apierrors"
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
					customError := apierrors.ERROR_CODE_UNIQUE_KEY_VIOLATION.NewApiError("A user with the same email already exists", result.Error)
					return customError
				}
				// Fallback for any other unique constraint violation not specifically handled
				customError := apierrors.ERROR_CODE_UNIQUE_KEY_VIOLATION.NewApiError("A user with similar details already exists", result.Error)
				return customError
			}
		}

		// Fallback for any other database errors (e.g., connection issues, other integrity errors)
		customError := apierrors.ERROR_CODE_INTERNAL_SERVER.NewApiError("An unexpected database error occurred creating the entity", result.Error)
		return customError
	}

	return nil
}
