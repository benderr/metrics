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

	"github.com/benderr/metrics/internal/gserver/grpcapp"
	"github.com/benderr/metrics/internal/server/app"
	"github.com/benderr/metrics/internal/server/config"
	"github.com/benderr/metrics/pkg/logger"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	l, sync := logger.New()
	defer sync()

	l.Infoln("Build version:", buildVersion)
	l.Infoln("Build date:", buildDate)
	l.Infoln("Build commit:", buildCommit)

	c := config.MustLoad()

	l.Infow(
		"Starting server with",
		"config", c,
	)

	ctx := context.Background()

	if c.Transport == "grpc" {
		if err := grpcapp.New(c, l).Run(ctx); err != nil {
			panic(err)
		}
	} else if c.Transport == "http" {
		if err := app.New(c, l).Run(ctx); err != nil {
			panic(err)
		}
	} else {
		panic("unsupported transport")
	}

}
