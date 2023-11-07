package repository

import (
	"context"
	"fmt"
	"strings"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) GetStringValue() string {
	switch m.MType {
	case "gauge":
		return strings.TrimRight(fmt.Sprintf("%.3f", *m.Value), "0")
	case "counter":
		return fmt.Sprintf("%v", *m.Delta)
	}
	return ""
}

type MetricRepository interface {
	Update(metric Metrics) (*Metrics, error)
	Get(id string) (*Metrics, error)
	GetList() ([]Metrics, error)
	PingContext(ctx context.Context) error
}
