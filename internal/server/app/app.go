// Package app return instance of App for starting server
package app

import (
	"context"
	"net/http"

	"net/http/pprof"

	"github.com/go-chi/chi"

	"github.com/benderr/metrics/internal/server/config"
	"github.com/benderr/metrics/internal/server/handlers"
	"github.com/benderr/metrics/internal/server/logger"
	"github.com/benderr/metrics/internal/server/middleware/gziper"
	"github.com/benderr/metrics/internal/server/middleware/mlogger"
	"github.com/benderr/metrics/internal/server/middleware/sign"
	"github.com/benderr/metrics/internal/server/repository/storage"
)

// App consisting only one method Run to start server
type App struct {
	config *config.Config
	log    logger.Logger
}

// New return a new App object
func New(config *config.Config, log logger.Logger) *App {
	return &App{
		config: config,
		log:    log,
	}
}

// Run create storage which depends on config and listens on the TCP network address addr and then calls
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

	// register pprof methods
	chiRouter.Route("/debug/pprof", func(r chi.Router) {
		r.HandleFunc("/cmdline", pprof.Cmdline)
		r.HandleFunc("/profile", pprof.Profile)
		r.HandleFunc("/symbol", pprof.Symbol)
		r.HandleFunc("/trace", pprof.Trace)
		r.HandleFunc("/*", pprof.Index)
	})

	h.AddHandlers(chiRouter)

	return http.ListenAndServe(string(a.config.Server), chiRouter)
}
