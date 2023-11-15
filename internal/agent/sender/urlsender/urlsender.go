package urlsender

import (
	"errors"
	"fmt"

	"github.com/benderr/metrics/internal/metrics/report"
	"github.com/go-resty/resty/v2"
)

// Первая версия с передачей данных внутри url
func New(server string) *URLSender {
	return &URLSender{
		client: resty.New().SetBaseURL(server),
	}
}

type URLSender struct {
	client *resty.Client
}

func (h *URLSender) Send(metrics []report.MetricItem) error {
	allErrors := make([]error, 0)
	for _, metric := range metrics {
		switch metric.MType {
		case "counter":
			url := fmt.Sprintf("/%v/%v/%v/%v", "update", "counter", metric.ID, *metric.Delta)
			_, err := h.client.R().Post(url)
			allErrors = append(allErrors, err)

		case "gauge":
			url := fmt.Sprintf("/%v/%v/%v/%v", "update", "gauge", metric.ID, *metric.Value)
			_, err := h.client.R().Post(url)
			allErrors = append(allErrors, err)
		}
	}

	return errors.Join(allErrors...)
}
