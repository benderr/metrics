package dbstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/benderr/metrics/internal/server/repository"
	"github.com/benderr/metrics/pkg/retry"
)

// Оборачиваем репозиторий с добавлением возможности повтора операции при ошибках
// В качестве примера сделал для трех операций
func NewWithRetry(db *sql.DB, log repository.Logger) *MetricDBWithRetryRepository {
	repo := &MetricDBRepository{
		db:  db,
		log: log,
	}
	return &MetricDBWithRetryRepository{
		MetricDBRepository: repo,
	}
}

type MetricDBWithRetryRepository struct {
	*MetricDBRepository
}

func (m *MetricDBWithRetryRepository) Update(ctx context.Context, mtr repository.Metrics) (*repository.Metrics, error) {
	return retry.DoWithValue[*repository.Metrics](func() (*repository.Metrics, error) {
		return m.MetricDBRepository.Update(ctx, mtr)
	}, canRetry)
}

func (m *MetricDBWithRetryRepository) BulkUpdate(ctx context.Context, metrics []repository.Metrics) error {
	return retry.Do(func() error {
		return m.MetricDBRepository.BulkUpdate(ctx, metrics)
	}, canRetry)
}

func (m *MetricDBWithRetryRepository) PingContext(ctx context.Context) error {
	return retry.Do(func() error {
		return m.MetricDBRepository.PingContext(ctx)
	}, canRetry)
}

// Ретрай с обработкой ошибок sql Class 08
func isPgErrors(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
		case pgerrcode.ConnectionException:
		case pgerrcode.ConnectionDoesNotExist:
		case pgerrcode.ConnectionFailure:
		case pgerrcode.SQLClientUnableToEstablishSQLConnection:
		case pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection:
			return true
		default:
			return false
		}
	}
	return false
}

func canRetry(attempt int, err error) bool {
	if isPgErrors(err) {
		return retry.DefaultRetryCondition(attempt, err)
	}
	return false
}
