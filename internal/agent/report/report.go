package report

import (
	"fmt"
	"sync"

	"github.com/benderr/metrics/internal/agent/stats"
)

type Report struct {
	MetricItems map[string]MetricItem
	mu          sync.Mutex
}

type MetricItem struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func New() *Report {
	return &Report{
		MetricItems: make(map[string]MetricItem),
	}
}

func (r *Report) updateCounter(name string, value int64) {
	var newVal int64 = value
	if v, ok := r.MetricItems[name]; ok {
		newVal = *v.Delta + value
	}
	r.MetricItems[name] = MetricItem{
		ID:    name,
		Delta: &newVal,
		MType: "counter",
	}
}

func (r *Report) updateGauge(name string, value float64) {
	r.MetricItems[name] = MetricItem{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}
}

func (r *Report) Update(items []stats.Item) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, item := range items {
		switch item.Type {
		case "gauge":
			r.updateGauge(item.Name, item.Value)
		case "counter":
			r.updateCounter(item.Name, item.Delta)
		}
	}
}

func (r *Report) GetList() []MetricItem {
	r.mu.Lock()
	defer r.mu.Unlock()
	metrics := make([]MetricItem, 0)

	for _, value := range r.MetricItems {
		metrics = append(metrics, value)
	}
	fmt.Println("GetList", metrics)
	return metrics
}
