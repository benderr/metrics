package main

import (
	"fmt"

	"github.com/benderr/metrics/cmd/metrics"
)

func main() {
	ParseConfig()

	fmt.Printf("Started with params: \n -address %v\n -report interval %v \n -pool interval %v \n\n", config.Server, config.ReportInterval, config.PollInterval)

	<-metrics.StartAgent(config.PollInterval, config.ReportInterval, string(config.Server))
}
