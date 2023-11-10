package main

import (
	"log"

	"github.com/benderr/metrics/internal/agentconfig"
	"github.com/benderr/metrics/internal/metrics/agent"
	"github.com/benderr/metrics/internal/metrics/sender/bulksender"
	"go.uber.org/zap"
)

func main() {
	config, err := agentconfig.Parse()

	if err != nil {
		log.Fatal(err)
	}

	l, lerr := zap.NewDevelopment()
	if lerr != nil {
		panic(lerr)
	}
	defer l.Sync()

	sugar := *l.Sugar()

	sugar.Infow(
		"Started with params",
		"-address", config.Server,
		"-report interval", config.ReportInterval,
		"-pool interval", config.PollInterval,
	)

	//sender := urlsender.New(string(config.Server))
	//sender := jsonsender.New(string(config.Server))
	sender := bulksender.New(string(config.Server), &sugar)
	a := agent.New(config.PollInterval, config.ReportInterval, sender)

	<-a.Run()
}
