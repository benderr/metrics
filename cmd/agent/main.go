package main

import (
	"fmt"
	"log"

	"github.com/benderr/metrics/cmd/config/agentconfig"
	"github.com/benderr/metrics/internal/metrics"
)

func main() {
	config, err := agentconfig.Parse()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Started with params: \n -address %v\n -report interval %v \n -pool interval %v \n\n", config.Server, config.ReportInterval, config.PollInterval)

	agent := metrics.NewAgent(config.PollInterval, config.ReportInterval, string(config.Server))

	<-agent.Run()
}
