package sender

import (
	"context"

	"github.com/benderr/metrics/internal/agent/report"
)

type MetricSender interface {
	Send(ctx context.Context, metrics []report.MetricItem) error
}
