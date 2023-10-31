package main

import (
	"net/http"

	"github.com/benderr/metrics/internal/dump"
	"github.com/benderr/metrics/internal/filedump"
	"github.com/benderr/metrics/internal/handlers"
	"github.com/benderr/metrics/internal/middleware/gziper"
	"github.com/benderr/metrics/internal/middleware/logger"
	"github.com/benderr/metrics/internal/repository"
	"github.com/benderr/metrics/internal/serverconfig"
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

	var repo repository.MetricRepository = storage.New()

	//configure dumper
	f := filedump.New(config.FileStoragePath) //создаем ReadWriteCloser, тут в перспективе может быть не только файл
	dumper := dump.New(repo, &sugar, f)

	if config.Restore {
		dumper.Restore()
	}

	if config.StoreInterval == 0 {
		//оборачиваем репозиторий чтобы ловить Update и делать синхронную запись в файл
		repo = dumper.TrackRepository(repo)
	} else {
		go dumper.SaveByTime(config.StoreInterval)
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
