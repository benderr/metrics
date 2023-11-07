package filestorage

import (
	"encoding/json"
	"io"

	"github.com/benderr/metrics/internal/repository"
	"github.com/benderr/metrics/internal/repository/inmemory"
)

type FileMetricRepository struct {
	sync bool
	repository.MetricRepository
	rwg    ReadWriteGetter
	logger ErrorLogger
}

type ErrorLogger interface {
	Errorln(args ...interface{})
}

type ReadWriteGetter interface {
	Get() (io.ReadWriteCloser, error)
}

func New(rwg ReadWriteGetter, logger ErrorLogger, sync bool, restore bool) *FileMetricRepository {
	var repo repository.MetricRepository = inmemory.New()

	newRepo := &FileMetricRepository{
		sync:             sync,
		MetricRepository: repo,
		logger:           logger,
		rwg:              rwg,
	}

	if restore {
		newRepo.Restore()
	}

	return newRepo
}

func (m *FileMetricRepository) Update(metric repository.Metrics) (*repository.Metrics, error) {
	res, err := m.MetricRepository.Update(metric)
	if err == nil && m.sync {
		m.Sync()
	}

	return res, err
}

func (m *FileMetricRepository) Sync() error {
	w, err := m.rwg.Get()
	if err != nil {
		m.logger.Errorln("invalid writer", err)
		return err
	}

	list, err := m.GetList()
	if err != nil {
		m.logger.Errorln("data error", err)
		return err
	}

	defer w.Close()

	encoder := json.NewEncoder(w)
	for _, item := range list {
		encoder.Encode(item)
	}
	return nil
}

func (m *FileMetricRepository) Restore() error {
	r, err := m.rwg.Get()
	if err != nil {
		m.logger.Errorln("invalid reader", err)
		return err
	}
	defer r.Close()

	decoder := json.NewDecoder(r)

	for {
		metric := &repository.Metrics{}
		err := decoder.Decode(&metric)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		m.Update(*metric)
	}
	return nil
}
