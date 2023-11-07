package dbstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/benderr/metrics/internal/repository"
)

func New(db *sql.DB) *MetricDBRepository {
	return &MetricDBRepository{
		db: db,
	}
}

type MetricDBRepository struct {
	db *sql.DB
}

func (m *MetricDBRepository) Update(mtr repository.Metrics) (*repository.Metrics, error) {
	return nil, errors.New("no implemented")
}

func (m *MetricDBRepository) Get(id string) (*repository.Metrics, error) {
	return nil, errors.New("no implemented")
}

func (m *MetricDBRepository) GetList() ([]repository.Metrics, error) {
	return nil, errors.New("no implemented")
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
