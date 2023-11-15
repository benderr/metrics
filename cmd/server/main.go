package main

import (
	"context"
	"net/http"

	"github.com/benderr/metrics/internal/server/config"
	"github.com/benderr/metrics/internal/server/handlers"
	"github.com/benderr/metrics/internal/server/logger"
	"github.com/benderr/metrics/internal/server/middleware/gziper"
	"github.com/benderr/metrics/internal/server/middleware/mlogger"
	"github.com/benderr/metrics/internal/server/repository/storage"
	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config, confError := config.Parse()
	if confError != nil {
		panic(confError)
	}

	//configure logger
	l, sync := logger.New()

	defer sync()

	l.Infow(
		"Starting server with",
		"config", config,
	)

	//configure repo
	ctx := context.Background()
	repo, close := storage.New(ctx, config, l)
	defer close()

	//configure api
	h := handlers.NewHandlers(repo, l)
	mlog := mlogger.New(l)
	gzip := gziper.New(1, "application/json", "text/html")

	chiRouter := chi.NewRouter()
	chiRouter.Use(mlog.Middleware)
	chiRouter.Use(gzip.TransformWriter)
	chiRouter.Use(gzip.TransformReader)
	h.AddHandlers(chiRouter)

	err := http.ListenAndServe(string(config.Server), chiRouter)
	if err != nil {
		panic(err)
	}
}
