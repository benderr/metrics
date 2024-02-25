// Package app return instance of App for starting server
package app

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/benderr/metrics/internal/server/config"
	"github.com/benderr/metrics/internal/server/handlers"
	"github.com/benderr/metrics/internal/server/logger"
	"github.com/benderr/metrics/internal/server/middleware/mlogger"
	"github.com/benderr/metrics/internal/server/middleware/sign"
	"github.com/benderr/metrics/internal/server/repository/storage"
	"github.com/benderr/metrics/pkg/gziper"
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

	chiRouter.Mount("/debug", middleware.Profiler())

	h.AddHandlers(chiRouter)

	srv := http.Server{Addr: string(a.config.Server), Handler: chiRouter}

	ctxStop, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer stop()

	idleConnsClosed := make(chan struct{})

	go func() {
		<-ctxStop.Done()

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatal("shutdown error ", err)
		}

		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServeTLS(a.config.PublicKey, a.config.CryptoKey); err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	return nil

}
