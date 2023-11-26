package jsonsender

import (
	"errors"

	"github.com/benderr/metrics/internal/agent/apiclient"
	"github.com/benderr/metrics/internal/agent/report"
)

// Вторая версия, с передачей данных через json-body
func New(client *apiclient.Client) *JSONSender {
	return &JSONSender{
		client: client,
	}
}

type JSONSender struct {
	client *apiclient.Client
}

func (h *JSONSender) Send(metrics []report.MetricItem) error {
	allErrors := make([]error, 0)

	for _, metric := range metrics {
		m := metric
		go func() {
			h.client.R().
				SetHeader("Content-Type", "application/json").
				SetHeader("Accept-Encoding", "gzip").
				SetBody(m).
				Post("/update")
		}()

		//allErrors = append(allErrors, err)
	}

	return errors.Join(allErrors...)
}
