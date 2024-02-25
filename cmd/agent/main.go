// Package agent start process to collect metrics from gopsutil and runtime.MemStats.
//
// You can configure the following parameters:
//
// -a - server url for store metrics
//
// -r - report send to server interval (seconds)
//
// -p - create report interval (seconds)
//
// -k - secret key for signing request body
//
// For more information use:
//
//	cmd/server/server --help
package main

import (
	"context"
	"log"
	"time"

	"github.com/benderr/metrics/internal/agent/agent"
	agentconfig "github.com/benderr/metrics/internal/agent/config"
	"github.com/benderr/metrics/internal/agent/logger"
	"github.com/benderr/metrics/internal/agent/metricsender"
	"github.com/benderr/metrics/internal/agent/report"
	"github.com/benderr/metrics/internal/agent/stats/allstats"
	"github.com/benderr/metrics/internal/agent/stats/memstats"
	"github.com/benderr/metrics/internal/agent/stats/psstats"
	"github.com/benderr/metrics/internal/agent/ticker"
)

func main() {
	config, err := agentconfig.Parse()

	if err != nil {
		log.Fatal(err)
	}

	l, sync := logger.New()

	defer sync()

	l.Infow(
		"Started with params",
		"-address", config.Server,
		"-report interval", config.ReportInterval,
		"-pool interval", config.PollInterval,
		"-key ", config.SecretKey,
		"-config", config.ConfigFile,
		"-crypto-key", config.CryptoKey,
	)

	sender := metricsender.MustLoad(metricsender.BULK, config, l)

	ctx := context.Background()

	a := agent.New(sender, report.New())

	stats1 := memstats.New(time.Second * time.Duration(config.PollInterval))
	stats2 := psstats.New(time.Second * time.Duration(config.PollInterval))

	statCh := allstats.Join(ctx, stats1, stats2)

	a.Run(ctx, statCh, ticker.New(ctx, time.Second*time.Duration(config.ReportInterval)))
}
