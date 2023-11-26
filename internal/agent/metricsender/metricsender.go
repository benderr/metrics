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

func MustLoad(mode SenderMode, config *agentconfig.EnvConfig, logger logger.Logger) sender.MetricSender {

	client := apiclient.New(string(config.Server), config.SecretKey, logger)
	client.SetCustomRetries(3)
	client.SetSignedHeader()

	var newsender sender.MetricSender

	switch mode {
	case URL:
		newsender = urlsender.New(client)
	case JSON:
		newsender = jsonsender.New(client)
	case BULK:
		newsender = bulksender.New(client, logger)
	default:
		panic("incorrect sender mode")
	}

	return newsender
}
