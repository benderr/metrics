package main

import (
	"fmt"
	"net/http"

	"github.com/benderr/metrics/cmd/handlers"
	"github.com/benderr/metrics/cmd/storage"
)

func main() {

	mux := http.NewServeMux()

	var store storage.MemoryRepository = &storage.MemStorage{}

	mux.HandleFunc(`/`, handlers.NotFoundHandler)
	mux.HandleFunc(`/update/counter/`, handlers.UpdateCounterMetricHandler(store))
	mux.HandleFunc(`/update/gauge/`, handlers.UpdateGaugeMetricHandler(store))

	fmt.Println("Starting")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
