package inmemory

import (
	"context"
	"sync"

	"github.com/benderr/metrics/internal/repository"
)

func New() *InMemoryMetricRepository {
	return &InMemoryMetricRepository{
		Metrics: make([]repository.Metrics, 0),
	}
}

type InMemoryMetricRepository struct {
	Metrics []repository.Metrics
	mu      sync.Mutex
}

func (m *InMemoryMetricRepository) Update(ctx context.Context, mtr repository.Metrics) (*repository.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	metric, err := m.Get(ctx, mtr.ID)
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

func (m *InMemoryMetricRepository) Get(ctx context.Context, id string) (*repository.Metrics, error) {
	for i, metric := range m.Metrics {
		if metric.ID == id {
			return &m.Metrics[i], nil
		}
	}
	return nil, nil
}

func (m *InMemoryMetricRepository) GetList(ctx context.Context) ([]repository.Metrics, error) {
	return m.Metrics, nil
}

func (m *InMemoryMetricRepository) PingContext(ctx context.Context) error {
	return nil
}

func (m *InMemoryMetricRepository) BulkUpdate(ctx context.Context, metrics []repository.Metrics) error {

	if len(metrics) == 0 {
		return nil
	}

	for _, v := range metrics {
		m.Update(ctx, v)
	}

	return nil
}
