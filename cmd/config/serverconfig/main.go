package serverconfig

import (
	"errors"
	"flag"
	"os"
	"regexp"
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
	Server ServerAddress
}

var config = Config{
	Server: ":8080",
}

func init() {
	flag.Var(&config.Server, "a", "address and port to run server")

}

func Parse() *Config {
	flag.Parse()

	if address, ok := os.LookupEnv("ADDRESS"); ok {
		config.Server.Set(address)
	}

	return &config
}
