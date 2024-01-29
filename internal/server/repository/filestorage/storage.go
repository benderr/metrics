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

// FileMetricRepository consisting of the core methods used by App
type FileMetricRepository struct {
	sync bool
	repository.MetricRepository
	filePath string
	logger   repository.Logger
}

// New returns a new FileMetricRepository object
// thatimplements the MetricRepository interface.
//
// This repository used an in-memory repository
// with the addition of additional methods for backup and restoring
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

// Update insert or update metric.
//
// If metric exist, then update delta and value field,
// otherwise new metric inserted
//
// If FileMetricRepository.sync=true then the metrics are also saved to the file
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

// BulkUpdate insert or update slice of metric.
//
// If FileMetricRepository.sync=true then the metrics are also saved to the file
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

// Sync saved metrics from memory to file
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

// Restore load metrics from file to memory
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
