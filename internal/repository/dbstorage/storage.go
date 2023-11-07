package dbstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/benderr/metrics/internal/repository"
)

func New(db *sql.DB) *MetricDbRepository {
	return &MetricDbRepository{
		db: db,
	}
}

type MetricDbRepository struct {
	db *sql.DB
}

func (m *MetricDbRepository) Update(mtr repository.Metrics) (*repository.Metrics, error) {
	return nil, errors.New("no implemented")
}

func (m *MetricDbRepository) Get(id string) (*repository.Metrics, error) {
	return nil, errors.New("no implemented")
}

func (m *MetricDbRepository) GetList() ([]repository.Metrics, error) {
	return nil, errors.New("no implemented")
}

func (m *MetricDbRepository) PingContext(ctx context.Context) error {
	if m.db == nil {
		return errors.New("no initialized")
	}

	if err := m.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
