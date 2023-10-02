package main

import (
	"net/http"

	"github.com/benderr/metrics/cmd/handlers"
	"github.com/benderr/metrics/cmd/storage"
)

func main() {

	var store storage.MemoryRepository = &storage.MemStorage{}

	r := handlers.MakeRouter(store)

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
