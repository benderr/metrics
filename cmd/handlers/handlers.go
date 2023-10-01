package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/benderr/metrics/cmd/storage"
	"github.com/benderr/metrics/cmd/validate"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func UpdateCounterMetricHandler(store storage.MemoryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		params := strings.Split(r.URL.Path, "/")
		if len(params) < 5 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		memType := params[2]
		name := params[3]
		value := params[4]

		if metric, err := validate.ParseCounter(memType, name, value); err == nil {

			store.UpdateCounter(metric)
			fmt.Println(store)
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}

func UpdateGaugeMetricHandler(store storage.MemoryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		params := strings.Split(r.URL.Path, "/")
		if len(params) < 5 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		memType := params[2]
		name := params[3]
		value := params[4]

		if metric, err := validate.ParseGauge(memType, name, value); err == nil {
			store.UpdateGauge(metric)
			fmt.Println(store)
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}
