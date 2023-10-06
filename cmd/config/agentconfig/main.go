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
}

var config = EnvConfig{
	Server:         "http://localhost:8080",
	ReportInterval: 10,
	PollInterval:   2,
}

func init() {
	flag.Var(&config.Server, "a", "address and port to run server (with http transport)")
	flag.IntVar(&config.ReportInterval, "r", 10, "report send to server interval (seconds)")
	flag.IntVar(&config.PollInterval, "p", 2, "create report interval (seconds)")
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
