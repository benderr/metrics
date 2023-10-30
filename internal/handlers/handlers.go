package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/benderr/metrics/internal/storage"
	"github.com/benderr/metrics/internal/validate"
	"github.com/go-chi/chi"
)

type AppHandlers struct {
	metricRepo MetricRepository
}

type MetricRepository interface {
	Update(metric storage.Metrics) (*storage.Metrics, error)
	Get(id string) (*storage.Metrics, error)
	GetList() ([]storage.Metrics, error)
}

func NewHandlers(repo MetricRepository) AppHandlers {
	return AppHandlers{
		metricRepo: repo,
	}
}

func (a *AppHandlers) AddHandlers(r *chi.Mux) {
	r.Get("/", a.GetMetricListHandler)
	r.Post("/update/{type}/{name}/{value}", a.UpdateMetricByURLHandler)
	r.Get("/value/{type}/{name}", a.GetMetricByURLHandler)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", a.UpdateMetricHandler)
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", a.GetMetricHandler)
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "invalid route "+r.RequestURI, http.StatusBadRequest)
		})
	})
}

func (a *AppHandlers) UpdateMetricByURLHandler(w http.ResponseWriter, r *http.Request) {
	memType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if metric, err := validate.ParseCounter(memType, name, value); err == nil {
		a.metricRepo.Update(*metric)
		w.WriteHeader(http.StatusOK)
		return
	}

	if metric, err := validate.ParseGauge(memType, name, value); err == nil {
		a.metricRepo.Update(*metric)
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (a *AppHandlers) GetMetricByURLHandler(w http.ResponseWriter, r *http.Request) {
	memType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	metric, err := a.metricRepo.Get(name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if metric == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if metric.MType != memType {
		http.Error(w, "invalid memType", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metric.GetStringValue()))
}

func (a *AppHandlers) GetMetricListHandler(w http.ResponseWriter, r *http.Request) {
	var output bytes.Buffer

	output.WriteString("<table>")

	metrics, err := a.metricRepo.GetList()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, counter := range metrics {
		fmt.Fprintf(&output, "<tr><td>%v</td><td>%v</td></tr>", counter.ID, counter.GetStringValue())
	}

	output.WriteString("<table>")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(output.Bytes())
}

func (a *AppHandlers) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var metric storage.Metrics

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMetric, err := a.metricRepo.Update(metric)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(&newMetric)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (a *AppHandlers) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var metric storage.Metrics

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exist, err := a.metricRepo.Get(metric.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exist == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	res, err := json.Marshal(&exist)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
