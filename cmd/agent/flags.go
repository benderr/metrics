package main

import (
	"errors"
	"flag"
	"regexp"
	"strings"
)

type ServerAddress string

func (address *ServerAddress) String() string {
	return string(*address)
}

func (address *ServerAddress) Set(flagValue string) error {
	if len(flagValue) == 0 {
		return errors.New("empty address not valid")
	}

	//resty нужен протокол, в тестах указывается без протокола, обходим
	if !strings.HasPrefix(flagValue, "https://") && !strings.HasPrefix(flagValue, "http://") {
		flagValue = "http://" + flagValue
	}

	reg := regexp.MustCompile(`^https?://[0-9A-Za-z\.]+(\:[0-9]+)?$`)

	if !reg.MatchString(flagValue) {
		return errors.New("invalid address and port")
	}
	*address = ServerAddress(flagValue)
	return nil
}

var server = new(ServerAddress)
var pollInterval int
var reportInterval int

func init() {
	*server = "http://localhost:8080"
	flag.Var(server, "a", "address and port to run server (with http transport)")
	flag.IntVar(&reportInterval, "r", 10, "report send to server interval (seconds)")
	flag.IntVar(&pollInterval, "p", 2, "create report interval (seconds)")
}
