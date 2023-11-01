package main

import (
	"fmt"
	"log"

	"github.com/benderr/metrics/internal/agentconfig"
	"github.com/benderr/metrics/internal/metrics/agent"
	sender "github.com/benderr/metrics/internal/metrics/sender"
)

func main() {
	config, err := agentconfig.Parse()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Started with params: \n -address %v\n -report interval %v \n -pool interval %v \n\n", config.Server, config.ReportInterval, config.PollInterval)

	var sender agent.MetricSender = sender.NewJSONSender(string(config.Server))
	a := agent.New(config.PollInterval, config.ReportInterval, sender)

	<-a.Run()
}
