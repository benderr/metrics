// Package handlers contain all endpoints with handlers to manage metrics
package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/benderr/metrics/internal/server/logger"
	"github.com/benderr/metrics/internal/server/repository"
	"github.com/benderr/metrics/pkg/sign"
)

// @Title MetricStorage API
// @Description Metrics manager
// @Version 1.0

// @Host localhost:8080

type AppHandlers struct {
	secret     string
	metricRepo repository.MetricRepository
	logger     logger.Logger
}

// metricsDto model info
// @Description metrics dto for fetch full information
type metricsDto struct {
	ID    string `json:"id"`   // unique metric name
	MType string `json:"type"` // metric type enum gauge или counter
}

// New returned object AppHandlers.
// Usage:
//
//	h := handlers.New(repo, logger, secret)
//	h.AddHandlers(chiRouter)
func New(repo repository.MetricRepository, logger logger.Logger, secret string) AppHandlers {
	return AppHandlers{
		metricRepo: repo,
		logger:     logger,
		secret:     secret,
	}
}

// AddHandlers registers handlers in *chi.Mux
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

// UpdateMetricByURLHandler handler to update metric.
//
// Information is received via URL.
func (a *AppHandlers) UpdateMetricByURLHandler(w http.ResponseWriter, r *http.Request) {
	memType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if metric, err := ParseCounter(memType, name, value); err == nil {
		a.metricRepo.Update(r.Context(), *metric)
		w.WriteHeader(http.StatusOK)
		return
	}

	if metric, err := ParseGauge(memType, name, value); err == nil {
		a.metricRepo.Update(r.Context(), *metric)
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

// GetMetricByURLHandler handler to get information about metric.
//
// Information is received via URL.
func (a *AppHandlers) GetMetricByURLHandler(w http.ResponseWriter, r *http.Request) {
	memType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	metric, err := a.metricRepo.Get(r.Context(), name)

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

// GetMetricListHandler handler for obtaining information about all metrics.
func (a *AppHandlers) GetMetricListHandler(w http.ResponseWriter, r *http.Request) {
	var output bytes.Buffer

	output.WriteString("<table>")

	metrics, err := a.metricRepo.GetList(r.Context())

	if err != nil {
		a.logger.Errorln(err)
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

// UpdateMetricHandler handler to update metric.
//
// Information is received from response.Body.
// @Description Create/update metric
// @Param metric body metricsDto true "metric ID and MType"
// @Success 200 {object} repository.Metrics
// @Failure 400 {string} string "Bad request, id not specified"
// @Failure 404 {string} string "Metric not found"
// @Failure 500 {string} string "Internal error"
// @Router /update [post]
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

	newMetric, err := a.metricRepo.Update(r.Context(), metric)

	if err != nil {
		a.logger.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(&newMetric)

	if err != nil {
		a.logger.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetMetricHandler handler to get information about metric.
//
// Information is received from response.Body.
// @Description Fetch metric info
// @Param metric body metricsDto true "metric ID and MType"
// @Success 200 {object} repository.Metrics
// @Failure 400 {string} string "Bad request, id not specified"
// @Failure 404 {string} string "Metric not found"
// @Failure 500 {string} string "Internal error"
// @Router /value [post]
func (a *AppHandlers) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var metric metricsDto

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exist, err := a.metricRepo.Get(r.Context(), metric.ID)

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
		a.logger.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if a.secret != "" {
		signhex := sign.New(a.secret, res)
		a.logger.Infoln("generated sign", signhex)
		w.Header().Set("HashSHA256", signhex)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (a *AppHandlers) PingDBHandler(w http.ResponseWriter, r *http.Request) {

	if err := a.metricRepo.PingContext(r.Context()); err != nil {
		http.Error(w, "could't connect to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// BulkUpdateHandler handler to update metrics,
// this method expected array of metrics in response.Body.
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

	err = a.metricRepo.BulkUpdate(r.Context(), metrics)

	if err != nil {
		a.logger.Infoln("internal error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
