package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/benderr/metrics/cmd/metrics"
	"github.com/go-resty/resty/v2"
)

var pollInterval = 2
var reportInterval = 10
var server = "http://localhost:8080"

func main() {
	client := resty.New()
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
			time.Sleep(time.Duration(pollInterval) * time.Second)
		}
	}()

	go func(m metrics.IMetrics) {
		for {
			time.Sleep(time.Duration(reportInterval) * time.Second)

			for name, value := range m.GetCounters() {
				go func(name string, value int) {
					url := fmt.Sprintf("%v/%v/%v/%v/%v", server, "update", "counter", name, value)
					fmt.Println("try send counter", url)
					response, err := client.R().Post(url)
					if err != nil {
						fmt.Println("error counter", err)
					} else {
						fmt.Println("finish counter", response.StatusCode())
					}
				}(name, value)
			}
			for name, value := range m.GetGauges() {
				go func(name string, value float64) {
					url := fmt.Sprintf("%v/%v/%v/%v/%v", server, "update", "gauge", name, value)
					fmt.Println("try send gauge", url)
					response, err := client.R().Post(url)
					if err != nil {
						fmt.Println("error counter", err)
					} else {
						fmt.Println("finish gauge", response.StatusCode())
					}
				}(name, value)
			}
		}
	}(stored)

	for v := range ch {
		fmt.Println("Refresh metrics", v)
		stored.Update()
	}

	wg.Wait()

}
