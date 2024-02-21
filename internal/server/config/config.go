package config

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
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

const (
	defaultStoreInterval int = 300
)

type Config struct {
	Server          ServerAddress `env:"ADDRESS"`
	StoreInterval   int           `env:"STORE_INTERVAL"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH"`
	Restore         bool          `env:"RESTORE"`
	DatabaseDsn     string        `env:"DATABASE_DSN"`
	SecretKey       string        `env:"KEY"`
	CryptoKey       string        `env:"CRYPTO_KEY"`
	ConfigFile      string        `env:"CONFIG"`
}

var config = Config{
	Server:          ":8080",
	StoreInterval:   defaultStoreInterval,
	FileStoragePath: "/tmp/metrics-db.json",
	Restore:         true,
	DatabaseDsn:     "",
	SecretKey:       "",
	CryptoKey:       "",
	ConfigFile:      "",
}

func init() {
	flag.Func("c", "file with json config", parseConfigFile)      // First try to parse config.json
	flag.Func("config", "file with json config", parseConfigFile) // First try to parse config.json
	flag.Var(&config.Server, "a", "address and port to run server")
	flag.IntVar(&config.StoreInterval, "i", defaultStoreInterval, "report save interval (seconds)")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/metrics-db.json", "report store file name")
	flag.BoolVar(&config.Restore, "r", true, "restore report from file")
	flag.StringVar(&config.DatabaseDsn, "d", "", "connection string for postgre")
	flag.StringVar(&config.SecretKey, "k", "", "sha256 based secret key")
	flag.StringVar(&config.CryptoKey, "crypto-key", "", "crypto file for TLS")
}

func MustLoad() *Config {
	flag.Parse()

	err := env.Parse(&config)

	if err != nil {
		panic(err)
	}

	return &config
}

type jsonConfig struct {
	Address         string `json:"address"`
	Restore         bool   `json:"restore"`
	StoreInterval   *int   `json:"store_interval"`
	FileStoragePath string `json:"store_file"`
	DatabaseDsn     string `json:"database_dsn"`
	CryptoKey       string `json:"crypto_key"`
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

	if fileConfig.StoreInterval != nil {
		config.StoreInterval = *fileConfig.StoreInterval
	}

	config.CryptoKey = fileConfig.CryptoKey
	config.Restore = fileConfig.Restore
	config.FileStoragePath = fileConfig.FileStoragePath
	config.DatabaseDsn = fileConfig.DatabaseDsn

	return nil
}
