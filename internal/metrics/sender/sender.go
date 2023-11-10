package sender

import "github.com/benderr/metrics/internal/metrics/report"

type MetricSender interface {
	Send(metrics []report.MetricItem) error
}

type Logger interface {
	Infoln(args ...interface{})
	Errorln(args ...interface{})
}
