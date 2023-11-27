package main

import (
	"context"
	"log"

	"github.com/benderr/metrics/internal/agent/agent"
	agentconfig "github.com/benderr/metrics/internal/agent/config"
	"github.com/benderr/metrics/internal/agent/logger"
	"github.com/benderr/metrics/internal/agent/metricsender"
	"github.com/benderr/metrics/internal/agent/stats/allstats"
	"github.com/benderr/metrics/internal/agent/stats/memstats"
	"github.com/benderr/metrics/internal/agent/stats/psstats"
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
	)

	sender := metricsender.MustLoad(metricsender.BULK, config, l)

	ctx := context.Background()

	a := agent.New(config.PollInterval, config.ReportInterval, sender)

	stats1 := memstats.New(config.PollInterval) //метрики из runtime
	stats2 := psstats.New(config.PollInterval)  //метрики из gopsutil

	statCh := allstats.Join(ctx, stats1, stats2)

	a.Run(ctx, statCh)
}
