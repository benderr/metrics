package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benderr/metrics/cmd/storage"
	"github.com/stretchr/testify/assert"
)

type MockMemoryStorage struct {
	Counters map[string]storage.MetricCounterInfo
	Gauges   map[string]storage.MetricGaugeInfo
}

func (m *MockMemoryStorage) UpdateCounter(counter storage.MetricCounterInfo) {
	if metric, ok := m.Counters[counter.Name]; ok {
		m.Counters[counter.Name] = storage.MetricCounterInfo{
			Value: metric.Value + counter.Value,
			Name:  metric.Name,
		}
	} else {
		m.Counters[counter.Name] = counter
	}
}

func (m *MockMemoryStorage) UpdateGauge(gauge storage.MetricGaugeInfo) {
	m.Gauges[gauge.Name] = gauge
}

func TestUpdateCounterMetricHandler(t *testing.T) {
	type want struct {
		code     int
		counters map[string]storage.MetricCounterInfo
	}
	tests := []struct {
		name         string
		url          string
		prevCounters map[string]storage.MetricCounterInfo
		want         want
	}{
		{
			url:  "/update/counter/test/2",
			name: "Add new counter metric",
			want: want{
				code: http.StatusOK,
				counters: map[string]storage.MetricCounterInfo{
					"test": {
						Name:  "test",
						Value: 2,
					},
				},
			},
		},
		{
			url:  "/update/counter/test/2",
			name: "Update exist counter metric",
			prevCounters: map[string]storage.MetricCounterInfo{
				"test": {
					Name:  "test",
					Value: 2,
				},
			},
			want: want{
				code: http.StatusOK,
				counters: map[string]storage.MetricCounterInfo{
					"test": {
						Name:  "test",
						Value: 4,
					},
				},
			},
		},
		{
			url:  "/update/counter/test/string",
			name: "Add metric with invalid data",
			want: want{
				code:     http.StatusBadRequest,
				counters: map[string]storage.MetricCounterInfo{},
			},
		},
		{
			url:  "/update/counter",
			name: "Add metric without params",
			want: want{
				code:     http.StatusNotFound,
				counters: map[string]storage.MetricCounterInfo{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			var store = MockMemoryStorage{
				Counters: make(map[string]storage.MetricCounterInfo),
				Gauges:   make(map[string]storage.MetricGaugeInfo),
			}

			if len(test.prevCounters) > 0 {
				store.Counters = test.prevCounters
			}

			UpdateCounterMetricHandler(&store)(w, request)

			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.EqualValues(t, store.Counters, test.want.counters)
		})
	}
}

func TestUpdateGaugeMetricHandler(t *testing.T) {
	type want struct {
		code   int
		gauges map[string]storage.MetricGaugeInfo
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			url:  "/update/gauge/test/2.0",
			name: "Add new gauge metric",
			want: want{
				code: http.StatusOK,
				gauges: map[string]storage.MetricGaugeInfo{
					"test": {
						Name:  "test",
						Value: 2.0,
					},
				},
			},
		},
		{
			url:  "/update/gauge/test/string",
			name: "Add metric with invalid data",
			want: want{
				code:   http.StatusBadRequest,
				gauges: map[string]storage.MetricGaugeInfo{},
			},
		},
		{
			url:  "/update/gauge",
			name: "Add metric without params",
			want: want{
				code:   http.StatusNotFound,
				gauges: map[string]storage.MetricGaugeInfo{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			var store = MockMemoryStorage{
				Counters: make(map[string]storage.MetricCounterInfo),
				Gauges:   make(map[string]storage.MetricGaugeInfo),
			}

			UpdateGaugeMetricHandler(&store)(w, request)

			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.EqualValues(t, store.Gauges, test.want.gauges)
		})
	}
}
