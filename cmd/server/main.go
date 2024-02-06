// Package server start server with endpoints for store metrics.
//
// Server collect metrics with two types: counter, gauge.
// Counter metric adds a value to an existing metric.
// Gauge metric overwrites existing value with new value.
//
// The server work in 3 mode (see config):
//
// In-memory mode (default mode). All metrics stored in-memory (key-value storage).
//
//	cmd/server/server
//
// Or filestorage ()
//
//	cmd/server/server -f "/tmp/example.json"
//
// Or database (postgresql)
//
//	cmd/server/server -d 'postgres://host:port/db'
//
// See other flags
//
//	cmd/server/server --help
package main

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/benderr/metrics/internal/server/app"
	"github.com/benderr/metrics/internal/server/config"
	"github.com/benderr/metrics/internal/server/logger"
)

func main() {
	config := config.MustLoad()

	logger, sync := logger.New()

	defer sync()

	logger.Infow(
		"Starting server with",
		"config", config,
	)

	app := app.New(config, logger)

	ctx := context.Background()
	err := app.Run(ctx)

	if err != nil {
		panic(err)
	}
}
