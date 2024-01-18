package filestorage

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/benderr/metrics/internal/retry"
	"github.com/benderr/metrics/internal/server/repository"
	"github.com/benderr/metrics/internal/server/repository/inmemory"
)

type FileMetricRepository struct {
	sync bool
	repository.MetricRepository
	filePath string
	logger   repository.Logger
}

func New(filePath string, sync bool, logger repository.Logger) *FileMetricRepository {
	var repo repository.MetricRepository = inmemory.New()

	return &FileMetricRepository{
		sync:             sync,
		MetricRepository: repo,
		logger:           logger,
		filePath:         filePath,
	}
}

func (f *FileMetricRepository) getFile() (io.ReadWriteCloser, error) {
	file, err := os.OpenFile(f.filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *FileMetricRepository) Update(ctx context.Context, metric repository.Metrics) (*repository.Metrics, error) {
	res, err := f.MetricRepository.Update(ctx, metric)
	if err != nil {
		return nil, err
	}

	if f.sync {
		f.Sync(ctx)
	}

	return res, err
}

func (f *FileMetricRepository) BulkUpdate(ctx context.Context, metrics []repository.Metrics) error {

	if len(metrics) == 0 {
		return nil
	}

	for _, v := range metrics {
		_, err := f.MetricRepository.Update(ctx, v)
		if err != nil {
			return err
		}
	}

	if f.sync {
		f.Sync(ctx)
	}

	return nil
}

func (f *FileMetricRepository) Sync(ctx context.Context) error {
	return retry.Do(func() error {
		w, err := f.getFile()

		if err != nil {
			f.logger.Errorln("invalid writer", err)
			return err
		}
		defer w.Close()

		list, err := f.GetList(ctx)
		if err != nil {
			f.logger.Errorln("data error", err)
			return err
		}

		encoder := json.NewEncoder(w)
		for _, item := range list {
			encoder.Encode(item)
		}
		return nil
	}, retry.DefaultRetryCondition)
}

func (f *FileMetricRepository) Restore(ctx context.Context) error {
	return retry.Do(func() error {
		r, err := f.getFile()
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
	}, retry.DefaultRetryCondition)
}
