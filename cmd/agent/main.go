package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/benderr/metrics/cmd/metrics"
	"github.com/go-resty/resty/v2"
)

func main() {
	ParseConfig()

	client := resty.New().SetBaseURL(string(config.Server))
	fmt.Printf("Started with params: \n -address %v\n -report interval %v \n -pool interval %v \n\n", config.Server, config.ReportInterval, config.PollInterval)
	ch := make(chan bool, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)

	stored := &metrics.Metrics{
		Gauges:  make(map[string]float64),
		Counter: make(map[string]int),
	}

	go func() {
		for {
			ch <- true
			time.Sleep(time.Duration(config.PollInterval) * time.Second)
		}
	}()

	go func(m metrics.MetricsReadWrite) {
		for {
			time.Sleep(time.Duration(config.ReportInterval) * time.Second)
			metrics.SendMetrics(m, client)
		}
	}(stored)

	for v := range ch {
		fmt.Println("Refresh metrics", v)
		stored.Write()
	}

	wg.Wait()

}
