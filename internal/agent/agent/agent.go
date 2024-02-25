package agent

import (
	"context"

	"github.com/benderr/metrics/internal/agent/report"
	"github.com/benderr/metrics/internal/agent/sender"
	"github.com/benderr/metrics/internal/agent/stats"
)

type Agent struct {
	sender sender.MetricSender
	report IReport
}

type IReport interface {
	Update(items []stats.Item)
	GetList() []report.MetricItem
}

func New(sender sender.MetricSender, report IReport) *Agent {
	return &Agent{
		sender: sender,
		report: report,
	}
}

func (a *Agent) SendMetrics(metrics []report.MetricItem) error {
	return a.sender.Send(metrics)
}

func (a *Agent) Run(ctx context.Context, in <-chan []stats.Item, sendSignal <-chan struct{}) {
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-in:
			a.report.Update(v)
		case <-sendSignal:
			a.SendMetrics(a.report.GetList())
		}
	}
}
