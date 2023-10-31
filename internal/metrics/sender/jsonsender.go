package sender

import (
	"fmt"

	"github.com/benderr/metrics/internal/metrics/report"
	"github.com/go-resty/resty/v2"
)

// Вторая версия, с передачей данных через json-body
func NewJSONSender(server string) *JSONSender {
	return &JSONSender{
		client: resty.New().SetBaseURL(server),
	}
}

type JSONSender struct {
	client *resty.Client
}

func (h *JSONSender) Send(metric *report.MetricItem) error {
	_, err := h.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(*metric).
		Post("/update")

	if err != nil {
		fmt.Println("send metric error", err)
	}
	return err
}
