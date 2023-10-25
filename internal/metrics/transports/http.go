package http

import (
	"errors"
	"fmt"

	"github.com/benderr/metrics/internal/metrics/report"
	"github.com/go-resty/resty/v2"
)

// Первая версия с передачей данных внутри url
func NewSimpleHttp(server string) *HttpSimpleTransport {
	return &HttpSimpleTransport{
		client: resty.New().SetBaseURL(server),
	}
}

type HttpSimpleTransport struct {
	client *resty.Client
}

func (h *HttpSimpleTransport) Send(metric *report.MetricItem) error {
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
func NewJsonHttp(server string) *HttpJsonTransport {
	return &HttpJsonTransport{
		client: resty.New().SetBaseURL(server),
	}
}

type HttpJsonTransport struct {
	client *resty.Client
}

func (h *HttpJsonTransport) Send(metric *report.MetricItem) error {
	_, err := h.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(*metric).
		Post("/update")

	if err != nil {
		fmt.Println("send metric error", err)
	}
	return err
}
