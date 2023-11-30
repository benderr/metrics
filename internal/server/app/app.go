package app

import (
	"context"
	"net/http"

	"github.com/benderr/metrics/internal/server/config"
	"github.com/benderr/metrics/internal/server/handlers"
	"github.com/benderr/metrics/internal/server/logger"
	"github.com/benderr/metrics/internal/server/middleware/gziper"
	"github.com/benderr/metrics/internal/server/middleware/mlogger"
	"github.com/benderr/metrics/internal/server/middleware/sign"
	"github.com/benderr/metrics/internal/server/repository/storage"
	"github.com/go-chi/chi"
)

type App struct {
	config *config.Config
	log    logger.Logger
}

func New(config *config.Config, log logger.Logger) *App {
	return &App{
		config: config,
		log:    log,
	}
}

func (a *App) Run(ctx context.Context) error {

	repo, err := storage.New(ctx, a.config, a.log)
	if err != nil {
		return err
	}

	h := handlers.New(repo, a.log, a.config.SecretKey)
	mwlog := mlogger.New(a.log)
	mwgzip := gziper.New(1, "application/json", "text/html")
	mwsign := sign.New(a.config.SecretKey, a.log)

	chiRouter := chi.NewRouter()
	chiRouter.Use(mwsign.CheckSign)
	chiRouter.Use(mwlog.Middleware)
	chiRouter.Use(mwgzip.TransformWriter)
	chiRouter.Use(mwgzip.TransformReader)

	h.AddHandlers(chiRouter)

	return http.ListenAndServe(string(a.config.Server), chiRouter)
}
