package agent

import (
	"time"

	"github.com/benderr/metrics/internal/agent/report"
	"github.com/benderr/metrics/internal/agent/sender"
)

type Agent struct {
	PoolInterval   int
	ReportInterval int
	sender         sender.MetricSender
}

func New(poolInterval int, reportInterval int, sender sender.MetricSender) *Agent {
	return &Agent{
		PoolInterval:   poolInterval,
		ReportInterval: reportInterval,
		sender:         sender,
	}
}

func (a *Agent) SendMetrics(r *report.Report) error {
	metrics := make([]report.MetricItem, 0)

	for name, value := range r.Counters {
		value := value
		metrics = append(metrics, report.MetricItem{ID: name, MType: "counter", Delta: &value})
	}

	for name, value := range r.Gauges {
		value := value
		metrics = append(metrics, report.MetricItem{ID: name, MType: "gauge", Value: &value})
	}

	return a.sender.Send(metrics)
}

func (a *Agent) Run() <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		pollTicker := time.NewTicker(time.Second * time.Duration(a.PoolInterval))
		reportTicker := time.NewTicker(time.Second * time.Duration(a.ReportInterval))

		defer pollTicker.Stop()
		defer reportTicker.Stop()

		report := report.New()

		for {
			select {
			case <-pollTicker.C:
				report.UpdateReport()
			case <-reportTicker.C:
				a.SendMetrics(report)
			}
		}
	}()

	return done
}
