package inmemory

import (
	"context"

	"github.com/benderr/metrics/internal/repository"
)

func New() *InMemoryMetricRepository {
	return &InMemoryMetricRepository{
		Metrics: make([]repository.Metrics, 0),
	}
}

type InMemoryMetricRepository struct {
	Metrics []repository.Metrics
}

func (m *InMemoryMetricRepository) Update(mtr repository.Metrics) (*repository.Metrics, error) {
	metric, err := m.Get(mtr.ID)
	if err != nil {
		return nil, err
	}
	if metric != nil {
		switch mtr.MType {
		case "gauge":
			metric.Value = mtr.Value
		case "counter":
			newVal := *metric.Delta + *mtr.Delta
			metric.Delta = &newVal
		}
		return metric, nil
	} else {
		m.Metrics = append(m.Metrics, mtr)
		return &mtr, nil
	}
}

func (m *InMemoryMetricRepository) Get(id string) (*repository.Metrics, error) {
	for i, metric := range m.Metrics {
		if metric.ID == id {
			return &m.Metrics[i], nil
		}
	}
	return nil, nil
}

func (m *InMemoryMetricRepository) GetList() ([]repository.Metrics, error) {
	return m.Metrics, nil
}

func (m *InMemoryMetricRepository) PingContext(ctx context.Context) error {
	return nil
}
