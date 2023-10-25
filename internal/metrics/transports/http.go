package http

import (
	"errors"
	"fmt"

	"github.com/benderr/metrics/internal/metrics/report"
	"github.com/go-resty/resty/v2"
)

// Первая версия с передачей данных внутри url
func NewSimpleHTTP(server string) *HTTPSimpleTransport {
	return &HTTPSimpleTransport{
		client: resty.New().SetBaseURL(server),
	}
}

type HTTPSimpleTransport struct {
	client *resty.Client
}

func (h *HTTPSimpleTransport) Send(metric *report.MetricItem) error {
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

// Вторая версия, с передачей данных через json-body
func NewJSONHTTP(server string) *HTTPJSONTransport {
	return &HTTPJSONTransport{
		client: resty.New().SetBaseURL(server),
	}
}

type HTTPJSONTransport struct {
	client *resty.Client
}

func (h *HTTPJSONTransport) Send(metric *report.MetricItem) error {
	_, err := h.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(*metric).
		Post("/update")

	if err != nil {
		fmt.Println("send metric error", err)
	}
	return err
}
