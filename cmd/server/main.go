package main

import (
	"context"

	"github.com/benderr/metrics/internal/server/app"
	"github.com/benderr/metrics/internal/server/config"
	"github.com/benderr/metrics/internal/server/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
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
