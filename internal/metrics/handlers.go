package metrics

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

func sendMetrics(m *Metrics, client *resty.Client) {
	for name, value := range m.Counters {
		go func(name string, value int) {
			url := fmt.Sprintf("/%v/%v/%v/%v", "update", "counter", name, value)
			client.R().Post(url)
		}(name, value)
	}

	for name, value := range m.Gauges {
		go func(name string, value float64) {
			url := fmt.Sprintf("/%v/%v/%v/%v", "update", "gauge", name, value)
			client.R().Post(url)
		}(name, value)
	}
}

type Agent struct {
	PoolInterval   int
	ReportInterval int
	Server         string
}

func NewAgent(poolInterval int, reportInterval int, server string) *Agent {
	return &Agent{
		PoolInterval:   poolInterval,
		ReportInterval: reportInterval,
		Server:         server,
	}
}

func (agent *Agent) Run() <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		pollTicker := time.NewTicker(time.Second * time.Duration(agent.PoolInterval))
		reportTicker := time.NewTicker(time.Second * time.Duration(agent.ReportInterval))

		defer pollTicker.Stop()
		defer reportTicker.Stop()

		report := &Metrics{
			Gauges:   make(map[string]float64),
			Counters: make(map[string]int),
		}

		client := resty.New().SetBaseURL(agent.Server)

		for {
			select {
			case <-pollTicker.C:
				report.UpdateReport()
			case <-reportTicker.C:
				sendMetrics(report, client)
			}
		}
	}()

	return done
}
