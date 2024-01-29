package inmemory

import (
	"context"
	"sync"

	"github.com/benderr/metrics/internal/server/repository"
)

type keyValueMetricRepository struct {
	Metrics map[string]*repository.Metrics
	mu      sync.Mutex
}

func NewFast() *keyValueMetricRepository {
	return &keyValueMetricRepository{
		Metrics: make(map[string]*repository.Metrics, 0),
	}
}

func (m *keyValueMetricRepository) Update(ctx context.Context, mtr repository.Metrics) (*repository.Metrics, error) {
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
		m.Metrics[mtr.ID] = &mtr
		return &mtr, nil
	}
}

func (m *keyValueMetricRepository) Get(ctx context.Context, id string) (*repository.Metrics, error) {
	if v, ok := m.Metrics[id]; ok {
		return v, nil
	}
	return nil, nil
}

func (m *keyValueMetricRepository) GetList(ctx context.Context) ([]repository.Metrics, error) {
	res := make([]repository.Metrics, 0)
	for _, val := range m.Metrics {
		res = append(res, *val)
	}
	return res, nil
}

func (m *keyValueMetricRepository) PingContext(ctx context.Context) error {
	return nil
}

func (m *keyValueMetricRepository) BulkUpdate(ctx context.Context, metrics []repository.Metrics) error {

	if len(metrics) == 0 {
		return nil
	}

	for _, v := range metrics {
		m.Update(ctx, v)
	}

	return nil
}
