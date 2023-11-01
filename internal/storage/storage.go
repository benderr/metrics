package storage

import (
	"fmt"
	"strings"
)

type MemType string

const (
	Gauge   MemType = "gauge"
	Counter MemType = "counter"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) GetStringValue() string {
	switch m.MType {
	case string(Gauge):
		return strings.TrimRight(fmt.Sprintf("%.3f", *m.Value), "0")
	case string(Counter):
		return fmt.Sprintf("%v", *m.Delta)
	}
	return ""
}

func New() *InMemoryMetricRepository {
	return &InMemoryMetricRepository{
		Metrics: make([]Metrics, 0),
	}
}

type InMemoryMetricRepository struct {
	Metrics []Metrics
}

func (m *InMemoryMetricRepository) Update(mtr Metrics) (*Metrics, error) {
	metric, err := m.Get(mtr.ID)
	if err != nil {
		return nil, err
	}
	if metric != nil {
		switch mtr.MType {
		case "gauge":
			metric.Value = mtr.Value
		case "counter":
			newVal := *metric.Delta + *mtr.Delta
			metric.Delta = &newVal
		}
		return metric, nil
	} else {
		m.Metrics = append(m.Metrics, mtr)
		return &mtr, nil
	}
}

func (m *InMemoryMetricRepository) Get(id string) (*Metrics, error) {
	for i, metric := range m.Metrics {
		if metric.ID == id {
			return &m.Metrics[i], nil
		}
	}
	return nil, nil
}

func (m *InMemoryMetricRepository) GetList() ([]Metrics, error) {
	return m.Metrics, nil
}
