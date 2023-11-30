package jsonsender

import (
	"errors"

	"github.com/benderr/metrics/internal/agent/apiclient"
	"github.com/benderr/metrics/internal/agent/report"
	"github.com/benderr/metrics/internal/agent/sender/worker"
)

// Вторая версия, с передачей данных через json-body
func New(client *apiclient.Client, rateLimit int) *JSONSender {
	return &JSONSender{
		client:    client,
		rateLimit: rateLimit,
	}
}

type JSONSender struct {
	client    *apiclient.Client
	rateLimit int
}

func (h *JSONSender) Send(metrics []report.MetricItem) error {
	allErrors := worker.Run(h.rateLimit, metrics, func(mi *report.MetricItem) error {
		_, e := h.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept-Encoding", "gzip").
			SetBody(mi).
			Post("/update")
		return e
	})
	return errors.Join(allErrors...)
}
