package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/benderr/metrics/cmd/handlers"
	"github.com/benderr/metrics/cmd/storage"
)

func main() {
	flag.Parse()
	var store storage.MemoryRepository = &storage.MemStorage{}

	r := handlers.MakeRouter(store)

	fmt.Println("Started on", server)
	err := http.ListenAndServe(string(*server), r)
	if err != nil {
		panic(err)
	}
}
