package config

import (
	"errors"
	"flag"
	"regexp"

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

	reg := regexp.MustCompile(`^([0-9A-Za-z\.]+)?(\:[0-9]+)?$`)

	if !reg.MatchString(flagValue) {
		return errors.New("invalid address and port")
	}

	*address = ServerAddress(flagValue)
	return nil
}

type Config struct {
	Server          ServerAddress `env:"ADDRESS"`
	StoreInterval   int           `env:"STORE_INTERVAL"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH"`
	Restore         bool          `env:"RESTORE"`
	DatabaseDsn     string        `env:"DATABASE_DSN"`
	SecretKey       string        `env:"KEY"`
}

var config = Config{
	Server:          ":8080",
	StoreInterval:   300,
	FileStoragePath: "/tmp/metrics-db.json",
	Restore:         true,
	DatabaseDsn:     "",
	SecretKey:       "",
}

func init() {
	flag.Var(&config.Server, "a", "address and port to run server")
	flag.IntVar(&config.StoreInterval, "i", 300, "report save interval (seconds)")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/metrics-db.json", "report store file name")
	flag.BoolVar(&config.Restore, "r", true, "restore report from file")
	flag.StringVar(&config.DatabaseDsn, "d", "", "connection string for postgre")
	flag.StringVar(&config.SecretKey, "k", "", "sha256 based secret key")
}

func MustLoad() *Config {
	flag.Parse()

	err := env.Parse(&config)

	if err != nil {
		panic(err)
	}

	return &config
}
