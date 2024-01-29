package inmemory

import (
	"context"
	"sync"

	"github.com/benderr/metrics/internal/server/repository"
)

type KeyValueMetricRepository struct {
	Metrics map[string]*repository.Metrics
	mu      sync.Mutex
}

// NewFast return in-memory repository instance.
//
// This implementation uses a map that is faster than the version based on slice.
//
// It's safe for concurrent use by multiple
// goroutines.
func NewFast() *KeyValueMetricRepository {
	return &KeyValueMetricRepository{
		Metrics: make(map[string]*repository.Metrics, 0),
	}
}

func (m *KeyValueMetricRepository) Update(ctx context.Context, mtr repository.Metrics) (*repository.Metrics, error) {
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

// Get returned information about metric by ID
func (m *KeyValueMetricRepository) Get(ctx context.Context, id string) (*repository.Metrics, error) {
	if v, ok := m.Metrics[id]; ok {
		return v, nil
	}
	return nil, nil
}

func (m *KeyValueMetricRepository) GetList(ctx context.Context) ([]repository.Metrics, error) {
	res := make([]repository.Metrics, 0)
	for _, val := range m.Metrics {
		res = append(res, *val)
	}
	return res, nil
}

func (m *KeyValueMetricRepository) PingContext(ctx context.Context) error {
	return nil
}

// BulkUpdate insert or update slice of metric to map-storage.
func (m *KeyValueMetricRepository) BulkUpdate(ctx context.Context, metrics []repository.Metrics) error {

	if len(metrics) == 0 {
		return nil
	}

	for _, v := range metrics {
		m.Update(ctx, v)
	}

	return nil
}
