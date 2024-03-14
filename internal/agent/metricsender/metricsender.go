package metricsender

import (
	"log"
	"os"

	"github.com/benderr/metrics/internal/agent/apiclient"
	agentconfig "github.com/benderr/metrics/internal/agent/config"
	"github.com/benderr/metrics/internal/agent/sender"
	"github.com/benderr/metrics/internal/agent/sender/bulksender"
	"github.com/benderr/metrics/internal/agent/sender/grpcsender"
	"github.com/benderr/metrics/internal/agent/sender/jsonsender"
	"github.com/benderr/metrics/internal/agent/sender/urlsender"
	"github.com/benderr/metrics/pkg/ipcheck"
	"github.com/benderr/metrics/pkg/logger"
)

type SenderMode int

const maxRetries int = 3

// MustLoad fabric method returned service for send metric to server.
//
// Send strategy depends SenderMode.
//
// url - send the metric using the POST method, the information is sent in the URL.
//
// json - send the metric using the POST method, the information is sent in the request.Body.
//
// bulk - send metrics in batches using the POST method.
//
// grpc_single - send one metric using the GRPS method.
//
// grpc_bulk - send metrics in batches using the GRPS method.
func MustLoad(config *agentconfig.EnvConfig, logger logger.Logger) (sender.MetricSender, func()) {

	noop := func() {}

	switch config.Mode {
	case "url":
		client := configureHttpClient(config, logger)
		newsender := urlsender.New(client, config.RateLimit)
		return newsender, noop
	case "json":
		client := configureHttpClient(config, logger)
		newsender := jsonsender.New(client, config.RateLimit)
		return newsender, noop
	case "bulk":
		client := configureHttpClient(config, logger)
		newsender := bulksender.New(client, logger)
		return newsender, noop
	case "grpc_single":
		return grpcsender.New(config.Server.GRPC(), logger, false)
	case "grpc_bulk":
		return grpcsender.New(config.Server.GRPC(), logger, true)
	default:
		log.Fatal("incorrect sender mode")
	}
	return nil, noop
}

func configureHttpClient(config *agentconfig.EnvConfig, logger logger.Logger) *apiclient.Client {
	client := apiclient.New(config.Server.HTTP(), config.SecretKey, logger)
	client.SetCustomRetries(maxRetries)
	client.SetSignedHeader()

	ip, err := ipcheck.GetHostIP()

	if err != nil {
		log.Fatal("can't set host IP", err)
	}
	logger.Infoln("HOST IP: ", ip.String())
	client.SetHeader("X-Real-IP", ip.String())

	if len(config.CryptoKey) > 0 {
		f, err := os.ReadFile(config.CryptoKey)
		if err != nil {
			log.Fatal("error open CryptoKey file", config.CryptoKey)
		}
		client.SetRootCertificateFromString(string(f))
		logger.Infoln("Certificate settled")
	}
	return client
}
