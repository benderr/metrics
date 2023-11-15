package bulksender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"time"

	"github.com/benderr/metrics/internal/agent/report"
	"github.com/benderr/metrics/internal/agent/sender"
	"github.com/go-resty/resty/v2"
)

// Третья версия, с передачей массива данных
func New(server string, log sender.Logger) *BulkSender {
	return &BulkSender{
		server: server,
		log:    log,
	}
}

type BulkSender struct {
	server string
	log    sender.Logger
}

func (b *BulkSender) Send(metrics []report.MetricItem) error {

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

	attempt := 0

	client := resty.
		New().
		SetBaseURL(b.server).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		SetRetryCount(3).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			attempt += 1
			wait := 0
			switch attempt {
			case 1:
				wait = 1
			case 2:
				wait = 3
			case 3:
				wait = 5
			}
			if wait > 0 {
				return time.Duration(wait) * time.Second, nil
			} else {
				return 0, errors.New("quota exceeded")
			}
		})

	res, err := client.
		R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetHeader("Content-Encoding", "gzip").
		SetBody(body).
		Post("/updates/")

	if err != nil {
		b.log.Errorln("bulk send error", err, string(res.Body()))
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
