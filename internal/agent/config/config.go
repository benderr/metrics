package agentconfig

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
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
	ConfigFile     string        `env:"CONFIG"`
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
	ConfigFile:     "",
}

func init() {
	flag.Func("c", "file with json config", parseConfigFile)      // First try to parse config.json
	flag.Func("config", "file with json config", parseConfigFile) // First try to parse config.json
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

type jsonConfig struct {
	Address        string `json:"address"`
	ReportInterval *int   `json:"report_interval"`
	PollInterval   *int   `json:"poll_interval"`
	CryptoKey      string `json:"crypto_key"`
}

func parseConfigFile(filePath string) error {
	config.ConfigFile = filePath

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	d := json.NewDecoder(f)
	fileConfig := &jsonConfig{}

	if err := d.Decode(fileConfig); err != nil {
		return err
	}

	if err = config.Server.Set(fileConfig.Address); err != nil {
		return err
	}

	if fileConfig.ReportInterval != nil {
		config.ReportInterval = *fileConfig.ReportInterval
	}

	if fileConfig.PollInterval != nil {
		config.PollInterval = *fileConfig.PollInterval
	}

	config.CryptoKey = fileConfig.CryptoKey

	return nil
}
