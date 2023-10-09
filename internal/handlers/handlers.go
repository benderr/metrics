package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/benderr/metrics/internal/storage"
	"github.com/benderr/metrics/internal/validate"
	"github.com/go-chi/chi"
)

type AppHandlers struct {
	store storage.MemoryRepository
}

func NewHandlers(store storage.MemoryRepository) AppHandlers {
	return AppHandlers{
		store: store,
	}
}

func (a *AppHandlers) NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", a.GetMetricListHandler())
	r.Post("/update/{type}/{name}/{value}", a.UpdateMetricHandler())
	r.Get("/value/{type}/{name}", a.GetMetricHandler())
	return r
}

func (a *AppHandlers) UpdateMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memType := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")

		if metric, err := validate.ParseCounter(memType, name, value); err == nil {
			a.store.UpdateCounter(metric)
			w.WriteHeader(http.StatusOK)
			return
		}

		if metric, err := validate.ParseGauge(memType, name, value); err == nil {
			a.store.UpdateGauge(metric)
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}

func (a *AppHandlers) GetMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memType := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")

		switch memType {
		case string(storage.Counter):
			if metric, ok := a.store.GetCounter(name); ok {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%v", metric.Value)))
				return
			}

		case string(storage.Gauge):
			if metric, ok := a.store.GetGauge(name); ok {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(formatGauge(metric.Value)))
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	}
}

func (a *AppHandlers) GetMetricListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var output bytes.Buffer

		output.WriteString("<table>")

		counters, err := a.store.GetCounters()
		gauges, err2 := a.store.GetGauges()

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

func formatGauge(v float64) string {
	//return strconv.FormatFloat(v, 'f', -1, 64)
	return strings.TrimRight(fmt.Sprintf("%.3f", v), "0")
}
