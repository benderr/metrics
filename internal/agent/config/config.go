package agentconfig

import (
	"errors"
	"flag"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v6"
)

type ServerAddress string

func (address *ServerAddress) String() string {
	return string(*address)
}

func (address *ServerAddress) Set(flagValue string) error {
	if len(flagValue) == 0 {
		return errors.New("empty address not valid")
	}

	reg := regexp.MustCompile(`^(https?://)?[0-9A-Za-z\.]+(\:[0-9]+)?$`)

	if !reg.MatchString(flagValue) {
		return errors.New("invalid address and port")
	}
	*address = ServerAddress(flagValue)
	return nil
}

type EnvConfig struct {
	Server         ServerAddress `env:"ADDRESS"`
	ReportInterval int           `env:"REPORT_INTERVAL"`
	PollInterval   int           `env:"POLL_INTERVAL"`
	SecretKey      string        `env:"KEY"`
	RateLimit      int           `env:"RATE_LIMIT"`
	CryptoKey      string        `env:"CRYPTO_KEY"`
}

const (
	defaultReportInterval int = 10
	defaultPoolInterval   int = 2
	defaultRateInterval   int = 10
)

var config = EnvConfig{
	Server:         "http://localhost:8080",
	ReportInterval: defaultReportInterval,
	PollInterval:   defaultPoolInterval,
	SecretKey:      "",
	RateLimit:      defaultRateInterval,
	CryptoKey:      "",
}

func init() {
	flag.Var(&config.Server, "a", "address and port to run server (with http transport)")
	flag.IntVar(&config.ReportInterval, "r", defaultReportInterval, "report send to server interval (seconds)")
	flag.IntVar(&config.PollInterval, "p", defaultPoolInterval, "create report interval (seconds)")
	flag.StringVar(&config.SecretKey, "k", "", "sha256 based secret key")
	flag.IntVar(&config.RateLimit, "l", defaultRateInterval, "rate limitter")
	flag.StringVar(&config.CryptoKey, "crypto-key", "", "crypto file for TLS")
}

func Parse() (*EnvConfig, error) {
	flag.Parse()

	err := env.Parse(&config)

	transformServerAddress(&config.Server)

	return &config, err
}

// resty нужен протокол, в тестах указывается без протокола, обходим
func transformServerAddress(address *ServerAddress) {
	if !strings.HasPrefix(address.String(), "https://") && !strings.HasPrefix(address.String(), "http://") {
		*address = ServerAddress("http://" + address.String())
	}
}
