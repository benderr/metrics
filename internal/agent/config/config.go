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

func (address *ServerAddress) HTTP() string {
	if !strings.HasPrefix(address.String(), "https://") && !strings.HasPrefix(address.String(), "http://") {
		return "http://" + address.String()
	}
	return address.String()
}

func (address *ServerAddress) GRPC() string {
	if strings.HasPrefix(address.String(), "https://") {
		return strings.TrimPrefix(address.String(), "https://")
	}

	if strings.HasPrefix(address.String(), "http://") {
		return strings.TrimPrefix(address.String(), "http://")
	}

	return address.String()
}

type EnvConfig struct {
	Server         ServerAddress `env:"ADDRESS"`
	ReportInterval int           `env:"REPORT_INTERVAL"`
	PollInterval   int           `env:"POLL_INTERVAL"`
	SecretKey      string        `env:"KEY"`
	RateLimit      int           `env:"RATE_LIMIT"`
	CryptoKey      string        `env:"CRYPTO_KEY"`
	ConfigFile     string        `env:"CONFIG"`
	Mode           string        `env:"MODE"`
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
	Mode:           "bulk",
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
	flag.StringVar(&config.Mode, "m", "bulk", "send to server mode: url, json, bulk, gsingle (grpc by metric), gbulk (grpc bulk flow), default bulk")
}

func Parse() (*EnvConfig, error) {
	flag.Parse()

	err := env.Parse(&config)

	return &config, err
}

type jsonConfig struct {
	Address        string `json:"address"`
	ReportInterval *int   `json:"report_interval"`
	PollInterval   *int   `json:"poll_interval"`
	CryptoKey      string `json:"crypto_key"`
	Mode           string `json:"mode"`
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

	config.Mode = fileConfig.Mode

	return nil
}
