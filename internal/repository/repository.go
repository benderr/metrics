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
		if m.Value != nil {
			return strings.TrimRight(fmt.Sprintf("%.3f", *m.Value), "0")
		} else {
			return "<nil>"
		}

	case "counter":
		if m.Delta != nil {
			return fmt.Sprintf("%v", *m.Delta)
		} else {
			return "<nil>"
		}
	}
	return ""
}

type MetricRepository interface {
	BulkUpdate(ctx context.Context, metrics []Metrics) error
	Update(ctx context.Context, metric Metrics) (*Metrics, error)
	Get(ctx context.Context, id string) (*Metrics, error)
	GetList(ctx context.Context) ([]Metrics, error)
	PingContext(ctx context.Context) error
}
