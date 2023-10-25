package main

import (
	"fmt"
	"log"

	"github.com/benderr/metrics/cmd/config/agentconfig"
	"github.com/benderr/metrics/internal/metrics/agent"
	http "github.com/benderr/metrics/internal/metrics/transports"
)

func main() {
	config, err := agentconfig.Parse()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Started with params: \n -address %v\n -report interval %v \n -pool interval %v \n\n", config.Server, config.ReportInterval, config.PollInterval)

	var sender agent.MetricSender = http.NewJSONHTTP(string(config.Server))
	a := agent.New(config.PollInterval, config.ReportInterval, sender)

	<-a.Run()
}
