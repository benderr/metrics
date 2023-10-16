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
	metricRepo MetricRepository
}

type MetricRepository interface {
	UpdateCounter(counter storage.MetricCounterInfo) error
	UpdateGauge(gauge storage.MetricGaugeInfo) error
	GetCounter(name string) (*storage.MetricCounterInfo, error)
	GetGauge(name string) (*storage.MetricGaugeInfo, error)
	GetCounters() ([]storage.MetricCounterInfo, error)
	GetGauges() ([]storage.MetricGaugeInfo, error)
}

func NewHandlers(repo MetricRepository) AppHandlers {
	return AppHandlers{
		metricRepo: repo,
	}
}

func (a *AppHandlers) NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", a.GetMetricListHandler)
	r.Post("/update/{type}/{name}/{value}", a.UpdateMetricHandler)
	r.Get("/value/{type}/{name}", a.GetMetricHandler)
	return r
}

func (a *AppHandlers) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	memType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if metric, err := validate.ParseCounter(memType, name, value); err == nil {
		a.metricRepo.UpdateCounter(metric)
		w.WriteHeader(http.StatusOK)
		return
	}

	if metric, err := validate.ParseGauge(memType, name, value); err == nil {
		a.metricRepo.UpdateGauge(metric)
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (a *AppHandlers) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	memType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	switch memType {
	case string(storage.Counter):
		if metric, _ := a.metricRepo.GetCounter(name); metric != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("%v", metric.Value)))
			return
		}

	case string(storage.Gauge):
		if metric, _ := a.metricRepo.GetGauge(name); metric != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(formatGauge(metric.Value)))
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (a *AppHandlers) GetMetricListHandler(w http.ResponseWriter, r *http.Request) {
	var output bytes.Buffer

	output.WriteString("<table>")

	counters, err := a.metricRepo.GetCounters()
	gauges, err2 := a.metricRepo.GetGauges()

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

func formatGauge(v float64) string {
	//return strconv.FormatFloat(v, 'f', -1, 64)
	return strings.TrimRight(fmt.Sprintf("%.3f", v), "0")
}
