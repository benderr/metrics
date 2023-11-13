package handlers

import (
	"errors"

	"github.com/benderr/metrics/internal/retry"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// Ретрай с обработкой ошибок sql Class 08
func CanRetry(attempt int, err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
		case pgerrcode.ConnectionException:
		case pgerrcode.ConnectionDoesNotExist:
		case pgerrcode.ConnectionFailure:
		case pgerrcode.SQLClientUnableToEstablishSQLConnection:
		case pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection:
			return retry.DefaultRetryCondition(attempt, err)
		default:
			return false
		}
		return false
	}
	return retry.DefaultRetryCondition(attempt, err)
}
