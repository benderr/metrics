package dump

import (
	"encoding/json"
	"io"
	"time"

	"github.com/benderr/metrics/internal/repository"
	"github.com/benderr/metrics/internal/storage"
)

type ErrorLogger interface {
	Errorln(args ...interface{})
}

type ReadWriteGetter interface {
	Get() (io.ReadWriteCloser, error)
}

type Dumper struct {
	metricRepo repository.MetricRepository
	logger     ErrorLogger
	rwg        ReadWriteGetter
}

func New(repo repository.MetricRepository, logger ErrorLogger, readWriteGetter ReadWriteGetter) *Dumper {
	return &Dumper{
		metricRepo: repo,
		logger:     logger,
		rwg:        readWriteGetter,
	}
}

func (m *Dumper) SaveByTime(storeIntervalSeconds int) {
	if storeIntervalSeconds == 0 {
		return
	}
	saveTicker := time.NewTicker(time.Second * time.Duration(storeIntervalSeconds))

	defer saveTicker.Stop()

	for {
		<-saveTicker.C
		go m.Save()
	}
}

func (m *Dumper) Save() error {
	w, err := m.rwg.Get()
	if err != nil {
		m.logger.Errorln("invalid writer", err)
		return err
	}

	list, err := m.metricRepo.GetList()
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

func (m *Dumper) Restore() error {
	r, err := m.rwg.Get()
	if err != nil {
		m.logger.Errorln("invalid reader", err)
		return err
	}
	defer r.Close()

	decoder := json.NewDecoder(r)

	for {
		metric := &storage.Metrics{}
		err := decoder.Decode(&metric)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		m.metricRepo.Update(*metric)
	}
	return nil
}

func (d *Dumper) TrackRepository(repo repository.MetricRepository) repository.MetricRepository {
	return &metricDumpRepository{repo, d}
}
