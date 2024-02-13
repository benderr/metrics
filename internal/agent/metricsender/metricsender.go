package metricsender

import (
	"github.com/benderr/metrics/internal/agent/apiclient"
	agentconfig "github.com/benderr/metrics/internal/agent/config"
	"github.com/benderr/metrics/internal/agent/logger"
	"github.com/benderr/metrics/internal/agent/sender"
	"github.com/benderr/metrics/internal/agent/sender/bulksender"
	"github.com/benderr/metrics/internal/agent/sender/jsonsender"
	"github.com/benderr/metrics/internal/agent/sender/urlsender"
)

type SenderMode int

const (
	JSON SenderMode = iota
	URL
	BULK
)

const MAX_RETRIES int = 3

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
	client.SetCustomRetries(MAX_RETRIES)
	client.SetSignedHeader()

	var newsender sender.MetricSender

	switch mode {
	case URL:
		newsender = urlsender.New(client, config.RateLimit)
	case JSON:
		newsender = jsonsender.New(client, config.RateLimit)
	case BULK:
		newsender = bulksender.New(client, logger)
	default:
		panic("incorrect sender mode")
	}

	return newsender
}
