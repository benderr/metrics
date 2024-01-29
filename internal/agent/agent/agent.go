package agent

import (
	"context"
	"time"

	"github.com/benderr/metrics/internal/agent/report"
	"github.com/benderr/metrics/internal/agent/sender"
	"github.com/benderr/metrics/internal/agent/stats"
)

type Agent struct {
	PoolInterval int
	sender       sender.MetricSender
}

func New(poolInterval int, sender sender.MetricSender) *Agent {
	return &Agent{
		PoolInterval: poolInterval,
		sender:       sender,
	}
}

func (a *Agent) SendMetrics(metrics []report.MetricItem) error {
	return a.sender.Send(metrics)
}

func (a *Agent) Run(ctx context.Context, in <-chan []stats.Item, sendSignal <-chan time.Time) {
	r := report.New()
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-in:
			r.Update(v)
		case <-sendSignal:
			a.SendMetrics(r.GetList())
		}
	}
}
