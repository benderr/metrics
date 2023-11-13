package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"os"

	"github.com/benderr/metrics/internal/repository"
)

func New(db *sql.DB, log ErrorLogger) *MetricDBRepository {
	return &MetricDBRepository{
		db:  db,
		log: log,
	}
}

type ErrorLogger interface {
	Errorln(args ...interface{})
}

type MetricDBRepository struct {
	db  *sql.DB
	log ErrorLogger
}

func (m *MetricDBRepository) Update(ctx context.Context, mtr repository.Metrics) (*repository.Metrics, error) {
	delta := sql.NullInt64{}
	value := sql.NullFloat64{}
	if mtr.Delta != nil {
		delta = sql.NullInt64{Valid: true, Int64: *mtr.Delta}
	}

	if mtr.Value != nil {
		value = sql.NullFloat64{Valid: true, Float64: *mtr.Value}
	}

	_, err := m.db.ExecContext(ctx, `INSERT INTO metrics (id, type, delta, value)
	VALUES($1, $2, $3, $4) 
	ON CONFLICT (id) 
	DO 
	   UPDATE SET delta=metrics.delta + $3, value=$4`, mtr.ID, mtr.MType, delta, value)

	if err != nil {
		return nil, err
	}

	return m.Get(ctx, mtr.ID)
}

func (m *MetricDBRepository) BulkUpdate(ctx context.Context, metrics []repository.Metrics) error {

	if len(metrics) == 0 {
		return nil
	}

	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	tx, err := m.db.BeginTx(newCtx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(newCtx, `INSERT INTO metrics (id, type, delta, value)
	VALUES($1, $2, $3, $4) 
	ON CONFLICT (id) 
	DO UPDATE SET delta=metrics.delta + $3, value=$4`)

	if err != nil {
		return err
	}

	for _, mtr := range metrics {
		delta := sql.NullInt64{}
		value := sql.NullFloat64{}
		if mtr.Delta != nil {
			delta = sql.NullInt64{Valid: true, Int64: *mtr.Delta}
		}

		if mtr.Value != nil {
			value = sql.NullFloat64{Valid: true, Float64: *mtr.Value}
		}
		_, err := stmt.ExecContext(newCtx, mtr.ID, mtr.MType, delta, value)

		if err != nil {
			stmt.Close()
			return err
		}
	}

	err = stmt.Close()

	if err != nil {
		return err
	}

	err = tx.Commit()

	if err != nil {
		return err
	}
	return err
}

func (m *MetricDBRepository) Get(ctx context.Context, id string) (*repository.Metrics, error) {
	row := m.db.QueryRowContext(ctx, "SELECT id, type, delta, value from metrics WHERE id = $1", id)
	var v repository.Metrics
	err := row.Scan(&v.ID, &v.MType, &v.Delta, &v.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &v, nil
}

func (m *MetricDBRepository) GetList(ctx context.Context) ([]repository.Metrics, error) {
	metrics := make([]repository.Metrics, 0)

	rows, err := m.db.QueryContext(ctx, "SELECT id, type, delta, value from metrics ORDER BY id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var v repository.Metrics
		err = rows.Scan(&v.ID, &v.MType, &v.Delta, &v.Value)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (m *MetricDBRepository) PingContext(ctx context.Context) error {
	if m.db == nil {
		return errors.New("no initialized")
	}

	if err := m.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (m *MetricDBRepository) Prepare(ctx context.Context) error {
	if err := m.PingContext(ctx); err != nil {
		return err
	}
	content, err := os.ReadFile("./init.sql")

	if err != nil {
		return err
	}

	sql := string(content)

	_, err = m.db.ExecContext(ctx, sql)

	if err != nil {
		return err
	}

	return nil
}
