package main

import (
	"net/http"

	"github.com/benderr/metrics/cmd/config/serverconfig"
	"github.com/benderr/metrics/internal/handlers"
	"github.com/benderr/metrics/internal/middleware"
	"github.com/benderr/metrics/internal/storage"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func main() {
	config := serverconfig.Parse()

	logger, logError := zap.NewDevelopment()
	if logError != nil {
		panic(logError)
	}
	defer logger.Sync()

	sugar = *logger.Sugar()

	sugar.Infow(
		"Starting server",
		"addr", config.Server,
	)

	var repo handlers.MetricRepository = storage.New()

	h := handlers.NewHandlers(repo)
	l := middleware.NewLoggingMiddleware(&sugar)
	chiRouter := chi.NewRouter()
	chiRouter.Use(l.Middleware)
	h.AddHandlers(chiRouter)

	err := http.ListenAndServe(string(config.Server), chiRouter)
	if err != nil {
		panic(err)
	}
}
