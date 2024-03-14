package urlsender

import (
	"context"
	"errors"
	"fmt"

	"github.com/benderr/metrics/internal/agent/apiclient"
	"github.com/benderr/metrics/internal/agent/report"
	"github.com/benderr/metrics/internal/agent/sender/worker"
)

// Первая версия с передачей данных внутри url
func New(client *apiclient.Client, rateLimit int) *URLSender {
	return &URLSender{
		client:    client,
		rateLimit: rateLimit,
	}
}

type URLSender struct {
	client    *apiclient.Client
	rateLimit int
}

func (h *URLSender) Send(ctx context.Context, metrics []report.MetricItem) error {
	allErrors := worker.Run(h.rateLimit, metrics, func(mi *report.MetricItem) error {
		switch mi.MType {
		case "counter":
			url := fmt.Sprintf("/%v/%v/%v/%v", "update", "counter", mi.ID, *mi.Delta)
			_, err := h.client.R().Post(url)
			return err

		case "gauge":
			url := fmt.Sprintf("/%v/%v/%v/%v", "update", "gauge", mi.ID, *mi.Value)
			_, err := h.client.R().SetContext(ctx).Post(url)
			return err
		}
		return nil
	})
	return errors.Join(allErrors...)
}
