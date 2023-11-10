package bulksender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"

	"github.com/benderr/metrics/internal/metrics/report"
	"github.com/benderr/metrics/internal/metrics/sender"
	"github.com/go-resty/resty/v2"
)

// Третья версия, с передачей массива данных
func New(server string, log sender.Logger) *BulkSender {
	return &BulkSender{
		client: resty.New().SetBaseURL(server),
		log:    log,
	}
}

type BulkSender struct {
	client *resty.Client
	log    sender.Logger
}

func (h *BulkSender) Send(metrics []report.MetricItem) error {

	if len(metrics) == 0 {
		return nil
	}

	bufBytes, err := json.Marshal(metrics)

	if err != nil {
		return err
	}

	body, err := compress(bufBytes)

	if err != nil {
		return err
	}

	res, err := h.client.R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetHeader("Content-Encoding", "gzip").
		SetBody(body).
		Post("/updates/")

	if err != nil {
		h.log.Errorln("bulk send error", err, string(res.Body()))
	}

	return err
}

func compress(s []byte) ([]byte, error) {
	var buf bytes.Buffer

	zipped := gzip.NewWriter(&buf)
	_, err := zipped.Write(s)
	if err != nil {
		return nil, err
	}

	zipped.Close()

	return buf.Bytes(), nil
}
