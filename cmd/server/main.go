package main

import (
	"io"
	"net/http"
	"os"

	"github.com/benderr/metrics/cmd/config/serverconfig"
	"github.com/benderr/metrics/internal/dump"
	"github.com/benderr/metrics/internal/handlers"
	"github.com/benderr/metrics/internal/middleware/gziper"
	"github.com/benderr/metrics/internal/middleware/logger"
	"github.com/benderr/metrics/internal/storage"
	"github.com/go-chi/chi"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func main() {
	config, confError := serverconfig.Parse()
	if confError != nil {
		panic(confError)
	}

	l, lerr := zap.NewDevelopment()
	if lerr != nil {
		panic(lerr)
	}
	defer l.Sync()

	sugar = *l.Sugar()

	sugar.Infow(
		"Starting server",
		"addr", config.Server,
	)

	var repo handlers.MetricRepository = storage.New()

	//configure dumper
	dumper := dump.New(repo, &sugar, getWriterFunc(config.FileStoragePath), getReaderFunc(config.FileStoragePath))

	if config.Restore {
		dumper.Restore()
	}

	if config.StoreInterval == 0 {
		//оборачиваем репозиторий чтобы ловить Update
		repo = &metricDumpRepository{repo, *dumper}
	} else {
		dumper.SaveByTime(config.StoreInterval)
	}

	//configure api
	h := handlers.NewHandlers(repo)
	log := logger.New(&sugar)
	gzip := gziper.New(1, "application/json", "text/html")

	chiRouter := chi.NewRouter()
	chiRouter.Use(log.Middleware)
	chiRouter.Use(gzip.TransformWriter)
	chiRouter.Use(gzip.TransformReader)
	h.AddHandlers(chiRouter)

	err := http.ListenAndServe(string(config.Server), chiRouter)
	if err != nil {
		panic(err)
	}
}

func getReaderFunc(filePath string) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) {
		file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

func getWriterFunc(filePath string) func() (io.WriteCloser, error) {
	return func() (io.WriteCloser, error) {
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

type metricDumpRepository struct {
	handlers.MetricRepository
	dumper dump.Dumper
}

func (m *metricDumpRepository) Update(metric storage.Metrics) (*storage.Metrics, error) {
	res, err := m.MetricRepository.Update(metric)
	if err == nil {
		m.dumper.Save()
	}

	return res, err
}
