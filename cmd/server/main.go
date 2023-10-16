package main

import (
	"fmt"
	"net/http"

	"github.com/benderr/metrics/cmd/config/serverconfig"
	"github.com/benderr/metrics/internal/handlers"
	"github.com/benderr/metrics/internal/storage"
)

func main() {
	config := serverconfig.Parse()

	var repo handlers.MetricRepository = &storage.InMemoryMetricRepository{
		Counters: make([]storage.MetricCounterInfo, 0),
		Gauges:   make([]storage.MetricGaugeInfo, 0),
	}

	h := handlers.NewHandlers(repo)

	fmt.Println("Started on", config.Server)

	err := http.ListenAndServe(string(config.Server), h.NewRouter())
	if err != nil {
		panic(err)
	}
}
