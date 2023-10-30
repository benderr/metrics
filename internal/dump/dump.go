package dump

import (
	"encoding/json"
	"io"
	"time"

	"github.com/benderr/metrics/internal/storage"
)

type MetricRepo interface {
	GetList() ([]storage.Metrics, error)
	Update(metric storage.Metrics) (*storage.Metrics, error)
}

type Logger interface {
	Errorln(args ...interface{})
}

type Dumper struct {
	metricRepo MetricRepo
	logger     Logger
	writer     func() (io.WriteCloser, error)
	reader     func() (io.ReadCloser, error)
}

func New(repo MetricRepo, logger Logger, writer func() (io.WriteCloser, error), reader func() (io.ReadCloser, error)) *Dumper {
	return &Dumper{
		metricRepo: repo,
		logger:     logger,
		writer:     writer,
		reader:     reader,
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

func (m *Dumper) Save() {
	w, err := m.writer()
	if err != nil {
		m.logger.Errorln("invalid writer", err)
		return
	}

	list, err := m.metricRepo.GetList()
	if err != nil {
		m.logger.Errorln("data error", err)
		return
	}

	defer w.Close()

	encoder := json.NewEncoder(w)
	for _, item := range list {
		encoder.Encode(item)
	}
}

func (m *Dumper) Restore() error {
	r, err := m.reader()
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
