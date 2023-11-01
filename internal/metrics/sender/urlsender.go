package sender

import (
	"errors"
	"fmt"

	"github.com/benderr/metrics/internal/metrics/report"
	"github.com/go-resty/resty/v2"
)

// Первая версия с передачей данных внутри url
func NewURLSender(server string) *URLSender {
	return &URLSender{
		client: resty.New().SetBaseURL(server),
	}
}

type URLSender struct {
	client *resty.Client
}

func (h *URLSender) Send(metric *report.MetricItem) error {
	switch metric.MType {
	case "counter":
		url := fmt.Sprintf("/%v/%v/%v/%v", "update", "counter", metric.ID, *metric.Delta)
		_, err := h.client.R().Post(url)
		if err != nil {
			fmt.Println("counter error", err)
		}
		return err
	case "gauge":
		url := fmt.Sprintf("/%v/%v/%v/%v", "update", "gauge", metric.ID, *metric.Value)
		_, err := h.client.R().Post(url)
		if err != nil {
			fmt.Println("gauge error", err)
		}
		return err
	}

	return errors.New("undefined MType")
}
