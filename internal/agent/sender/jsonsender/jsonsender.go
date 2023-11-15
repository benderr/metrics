package jsonsender

import (
	"errors"

	"github.com/benderr/metrics/internal/agent/report"
	"github.com/go-resty/resty/v2"
)

// Вторая версия, с передачей данных через json-body
func New(server string) *JSONSender {
	return &JSONSender{
		client: resty.New().SetBaseURL(server),
	}
}

type JSONSender struct {
	client *resty.Client
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
