package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/benderr/metrics/cmd/storage"
	"github.com/benderr/metrics/cmd/validate"
	"github.com/go-chi/chi"
)

func ListOfMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Список"))
}

func MakeRouter(store storage.MemoryRepository) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", GetMetricListHandler(store))
	r.Post("/update/{type}/{name}/{value}", UpdateMetricHandler(store))
	r.Get("/value/{type}/{name}", GetMetricHandler(store))
	return r
}

func UpdateMetricHandler(store storage.MemoryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memType := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")

		if metric, err := validate.ParseCounter(memType, name, value); err == nil {
			store.UpdateCounter(metric)
			w.WriteHeader(http.StatusOK)
			return
		}

		if metric, err := validate.ParseGauge(memType, name, value); err == nil {
			store.UpdateGauge(metric)
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetMetricHandler(store storage.MemoryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memType := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")

		switch memType {
		case string(storage.Counter):
			if metric, ok := store.GetCounter(name); ok {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%v", metric.Value)))
				return
			}

		case string(storage.Gauge):
			if metric, ok := store.GetGauge(name); ok {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%.3f", metric.Value)))
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	}
}

func GetMetricListHandler(store storage.MemoryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var output bytes.Buffer

		output.WriteString("<table>")

		counters, err := store.GetCounters()
		gauges, err2 := store.GetGauges()

		if err != nil || err2 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, counter := range counters {
			fmt.Fprintf(&output, "<tr><td>%v</td><td>%v</td></tr>", counter.Name, counter.Value)
		}

		for _, gauge := range gauges {
			fmt.Fprintf(&output, "<tr><td>%v</td><td>%f</td></tr>", gauge.Name, gauge.Value)
		}

		output.WriteString("<table>")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(output.Bytes())
	}
}
