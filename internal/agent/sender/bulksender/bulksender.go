package bulksender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"

	"github.com/benderr/metrics/internal/agent/apiclient"
	"github.com/benderr/metrics/internal/agent/report"
	"github.com/benderr/metrics/pkg/logger"
)

// Третья версия, с передачей массива данных
// Здесь нет реализации Worker Pool так как уходит всего один запрос со всеми метриками
func New(client *apiclient.Client, log logger.Logger) *BulkSender {
	return &BulkSender{
		client: client,
		log:    log,
	}
}

type BulkSender struct {
	client *apiclient.Client
	log    logger.Logger
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

	req := b.client.
		R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetHeader("Content-Encoding", "gzip").
		SetBody(body)

	_, err = req.
		Post("/updates/")

	if err != nil {
		b.log.Errorln("bulk send error", err)
	} else {
		b.log.Infoln("sent", string(bufBytes))
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
