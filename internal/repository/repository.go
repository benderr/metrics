package repository

import "github.com/benderr/metrics/internal/storage"

type MetricRepository interface {
	Update(metric storage.Metrics) (*storage.Metrics, error)
	Get(id string) (*storage.Metrics, error)
	GetList() ([]storage.Metrics, error)
}
