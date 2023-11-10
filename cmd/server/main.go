package main

import (
	"context"
	"net/http"

	"database/sql"

	"github.com/benderr/metrics/internal/dump"
	"github.com/benderr/metrics/internal/filedump"
	"github.com/benderr/metrics/internal/handlers"
	"github.com/benderr/metrics/internal/middleware/gziper"
	"github.com/benderr/metrics/internal/middleware/logger"
	"github.com/benderr/metrics/internal/repository"
	"github.com/benderr/metrics/internal/repository/dbstorage"
	"github.com/benderr/metrics/internal/repository/filestorage"
	"github.com/benderr/metrics/internal/repository/inmemory"
	"github.com/benderr/metrics/internal/serverconfig"
	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func main() {
	config, confError := serverconfig.Parse()
	if confError != nil {
		panic(confError)
	}

	//configure logger
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

	//configure repo
	var repo repository.MetricRepository
	var ctx = context.Background()

	switch {
	case config.DatabaseDsn != "":
		db, dberr := sql.Open("pgx", config.DatabaseDsn)
		if dberr != nil {
			panic(dberr)
		}
		defer db.Close()
		dbRepo := dbstorage.New(db)
		if err := dbRepo.Prepare(ctx); err != nil {
			panic(err)
		}
		repo = dbRepo

	case config.FileStoragePath != "":
		readWriter := filedump.New(config.FileStoragePath)
		sync := config.StoreInterval == 0
		fs := filestorage.New(readWriter, &sugar, sync)
		dumper := dump.New(fs)
		if !sync {
			go dumper.Start(ctx, config.StoreInterval)
		}
		if config.Restore {
			fs.Restore(ctx)
		}

		repo = fs

	default:
		repo = inmemory.New()
	}

	//configure api
	h := handlers.NewHandlers(repo, &sugar)
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
