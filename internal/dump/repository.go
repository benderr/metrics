package dump

import (
	"github.com/benderr/metrics/internal/repository"
	"github.com/benderr/metrics/internal/storage"
)

type metricDumpRepository struct {
	repository.MetricRepository
	dumper *Dumper
}

func (m *metricDumpRepository) Update(metric storage.Metrics) (*storage.Metrics, error) {
	res, err := m.MetricRepository.Update(metric)
	if err == nil {
		m.dumper.Save()
	}

	return res, err
}
