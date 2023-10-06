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
	var store storage.MemoryRepository = &storage.MemStorage{}

	r := handlers.NewRouter(store)

	fmt.Println("Started on", config.Server)
	err := http.ListenAndServe(string(config.Server), r)
	if err != nil {
		panic(err)
	}
}
