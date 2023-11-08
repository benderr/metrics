package filestorage

import (
	"context"
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

func New(rwg ReadWriteGetter, logger ErrorLogger, sync bool) *FileMetricRepository {
	var repo repository.MetricRepository = inmemory.New()

	return &FileMetricRepository{
		sync:             sync,
		MetricRepository: repo,
		logger:           logger,
		rwg:              rwg,
	}
}

func (f *FileMetricRepository) Update(ctx context.Context, metric repository.Metrics) (*repository.Metrics, error) {
	res, err := f.MetricRepository.Update(ctx, metric)
	if err == nil && f.sync {
		f.Sync(ctx)
	}

	return res, err
}

func (f *FileMetricRepository) Sync(ctx context.Context) error {
	w, err := f.rwg.Get()
	if err != nil {
		f.logger.Errorln("invalid writer", err)
		return err
	}

	list, err := f.GetList(ctx)
	if err != nil {
		f.logger.Errorln("data error", err)
		return err
	}

	defer w.Close()

	encoder := json.NewEncoder(w)
	for _, item := range list {
		encoder.Encode(item)
	}
	return nil
}

func (f *FileMetricRepository) Restore(ctx context.Context) error {
	r, err := f.rwg.Get()
	if err != nil {
		f.logger.Errorln("invalid reader", err)
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
		f.Update(ctx, *metric)
	}
	return nil
}
