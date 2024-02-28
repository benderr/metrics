package metricsender

import (
	"log"
	"os"

	"github.com/benderr/metrics/internal/agent/apiclient"
	agentconfig "github.com/benderr/metrics/internal/agent/config"
	"github.com/benderr/metrics/internal/agent/sender"
	"github.com/benderr/metrics/internal/agent/sender/bulksender"
	"github.com/benderr/metrics/internal/agent/sender/jsonsender"
	"github.com/benderr/metrics/internal/agent/sender/urlsender"
	"github.com/benderr/metrics/pkg/logger"
)

type SenderMode int

const (
	JSON SenderMode = iota
	URL
	BULK
)

const maxRetries int = 3

// MustLoad fabric method returned service for send metric to server.
//
// Send strategy depends SenderMode.
//
// URL - send the metric using the POST method, the information is sent in the URL.
//
// JSON - send the metric using the POST method, the information is sent in the request.Body.
//
// BULK - send metrics in batches using the POST method.
func MustLoad(mode SenderMode, config *agentconfig.EnvConfig, logger logger.Logger) sender.MetricSender {

	client := apiclient.New(string(config.Server), config.SecretKey, logger)
	client.SetCustomRetries(maxRetries)
	client.SetSignedHeader()

	if len(config.CryptoKey) > 0 {
		f, err := os.ReadFile(config.CryptoKey)
		if err != nil {
			log.Fatal("error open CryptoKey file", config.CryptoKey)
		}
		client.SetRootCertificateFromString(string(f))
		logger.Infoln("Certificate settled")
	}

	var newsender sender.MetricSender

	switch mode {
	case URL:
		newsender = urlsender.New(client, config.RateLimit)
	case JSON:
		newsender = jsonsender.New(client, config.RateLimit)
	case BULK:
		newsender = bulksender.New(client, logger)
	default:
		log.Fatal("incorrect sender mode")
	}

	return newsender
}
