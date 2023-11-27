package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/benderr/metrics/internal/agent/report"
	"github.com/benderr/metrics/internal/agent/sender"
	"github.com/benderr/metrics/internal/agent/stats"
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

func (a *Agent) SendMetrics(metrics []report.MetricItem) error {
	return a.sender.Send(metrics)
}

func (a *Agent) Run(ctx context.Context, in <-chan []stats.Item) {
	r := report.New()
	reportTicker := time.NewTicker(time.Second * time.Duration(a.ReportInterval))

	for {
		select {
		case <-ctx.Done():
			return
		case v := <-in:
			fmt.Println("before", v)
			r.Update(v)
			fmt.Println("after", v)
		case <-reportTicker.C:
			a.SendMetrics(r.GetList())
		}
	}
}
