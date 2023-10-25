package agent

import (
	"time"

	"github.com/benderr/metrics/internal/metrics/report"
)

type Agent struct {
	PoolInterval   int
	ReportInterval int
	sender         MetricSender
}

type MetricSender interface {
	Send(metric *report.MetricItem) error
}

func New(poolInterval int, reportInterval int, sender MetricSender) *Agent {
	return &Agent{
		PoolInterval:   poolInterval,
		ReportInterval: reportInterval,
		sender:         sender,
	}
}

func (a *Agent) SendMetrics(r *report.Report) {
	for name, value := range r.Counters {
		value := value
		m := &report.MetricItem{ID: name, MType: "counter", Delta: &value}
		go a.sender.Send(m)
	}

	for name, value := range r.Gauges {
		value := value
		m := &report.MetricItem{ID: name, MType: "gauge", Value: &value}
		go a.sender.Send(m)
	}
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
