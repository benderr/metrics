package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/benderr/metrics/internal/server/logger"
	"github.com/benderr/metrics/internal/server/repository"
	"github.com/go-chi/chi"
)

type AppHandlers struct {
	metricRepo repository.MetricRepository
	logger     logger.Logger
}

func NewHandlers(repo repository.MetricRepository, logger logger.Logger) AppHandlers {
	return AppHandlers{
		metricRepo: repo,
		logger:     logger,
	}
}

func (a *AppHandlers) AddHandlers(r *chi.Mux) {
	r.Get("/", a.GetMetricListHandler)
	r.Post("/update/{type}/{name}/{value}", a.UpdateMetricByURLHandler)
	r.Get("/value/{type}/{name}", a.GetMetricByURLHandler)
	r.Get("/ping", a.PingDBHandler)
	r.Post("/updates/", a.BulkUpdateHandler)

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

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	if metric, err := ParseCounter(memType, name, value); err == nil {
		a.metricRepo.Update(ctx, *metric)
		w.WriteHeader(http.StatusOK)
		return
	}

	if metric, err := ParseGauge(memType, name, value); err == nil {
		a.metricRepo.Update(ctx, *metric)
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (a *AppHandlers) GetMetricByURLHandler(w http.ResponseWriter, r *http.Request) {
	memType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	metric, err := a.metricRepo.Get(ctx, name)

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

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	metrics, err := a.metricRepo.GetList(ctx)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, counter := range metrics {
		fmt.Fprintf(&output, "<tr><td>%v</td><td>%s</td></tr>", counter.ID, counter.GetStringValue())
	}

	output.WriteString("<table>")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(output.Bytes())
}

func (a *AppHandlers) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var metric repository.Metrics

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	newMetric, err := a.metricRepo.Update(ctx, metric)

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
	var metric repository.Metrics

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	exist, err := a.metricRepo.Get(ctx, metric.ID)

	if err != nil {
		a.logger.Errorln("internal error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exist == nil {
		http.Error(w, "not found:", http.StatusNotFound)
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

func (a *AppHandlers) PingDBHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	if err := a.metricRepo.PingContext(ctx); err != nil {
		http.Error(w, "could't connect to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *AppHandlers) BulkUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	metrics := make([]repository.Metrics, 0)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		a.logger.Infoln("bad request:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		a.logger.Infoln("bad request unmarshal:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	err = a.metricRepo.BulkUpdate(ctx, metrics)

	if err != nil {
		a.logger.Infoln("internal error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
