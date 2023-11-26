package main

import (
	"log"

	"github.com/benderr/metrics/internal/agent/agent"
	agentconfig "github.com/benderr/metrics/internal/agent/config"
	"github.com/benderr/metrics/internal/agent/logger"
	"github.com/benderr/metrics/internal/agent/metricsender"
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

	a := agent.New(config.PollInterval, config.ReportInterval, sender)

	<-a.Run()
}
