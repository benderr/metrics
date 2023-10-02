package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benderr/metrics/cmd/storage"
	"github.com/go-resty/resty/v2"
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

func (m *MockMemoryStorage) GetCounters() ([]storage.MetricCounterInfo, error) {
	res := []storage.MetricCounterInfo{}
	for _, item := range m.Counters {
		res = append(res, item)
	}
	return res, nil
}

func (m *MockMemoryStorage) GetGauges() ([]storage.MetricGaugeInfo, error) {
	res := []storage.MetricGaugeInfo{}
	for _, item := range m.Gauges {
		res = append(res, item)
	}
	return res, nil
}

func (m *MockMemoryStorage) GetCounter(name string) (*storage.MetricCounterInfo, bool) {
	res, ok := m.Counters[name]
	return &res, ok
}

func (m *MockMemoryStorage) GetGauge(name string) (*storage.MetricGaugeInfo, bool) {
	res, ok := m.Gauges[name]
	return &res, ok
}

func TestUpdateMetricHandler(t *testing.T) {

	var store = MockMemoryStorage{
		Counters: make(map[string]storage.MetricCounterInfo),
		Gauges:   make(map[string]storage.MetricGaugeInfo),
	}

	r := MakeRouter(&store)
	server := httptest.NewServer(r)

	defer server.Close()

	type want struct {
		code int
	}
	tests := []struct {
		name   string
		url    string
		method string
		want   want
	}{
		{
			url:    "/update/counter/test/2",
			method: http.MethodPost,
			name:   "Add new counter metric",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			url:    "/update/counter/test/string",
			method: http.MethodPost,
			name:   "Add counter metric with invalid data",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			url:    "/update/counter",
			method: http.MethodPost,
			name:   "Add metric without params",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			url:    "/update/gauge/test/2.0",
			method: http.MethodPost,
			name:   "Add new gauge metric",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			url:    "/update/gauge/test/string",
			method: http.MethodPost,
			name:   "Add gauge metric with invalid data",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			url:    "/update/gauge",
			method: http.MethodPost,
			name:   "Add gauge metric without params",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			url:    "/update/gauge/test/2.0",
			method: http.MethodGet,
			name:   "Try Get Method",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = test.method
			req.URL = server.URL + test.url

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, test.want.code, resp.StatusCode())
		})
	}
}

func TestGetMetric(t *testing.T) {

	var store = MockMemoryStorage{
		Counters: map[string]storage.MetricCounterInfo{"test": {Name: "test", Value: 1}},
		Gauges:   map[string]storage.MetricGaugeInfo{"test2": {Name: "test2", Value: 100.123}},
	}

	r := MakeRouter(&store)
	server := httptest.NewServer(r)

	defer server.Close()

	type want struct {
		code    int
		content string
	}
	tests := []struct {
		name   string
		url    string
		method string
		want   want
	}{
		{
			url:    "/value/counter/test",
			method: http.MethodGet,
			name:   "Get counter test",
			want: want{
				code:    http.StatusOK,
				content: "1",
			},
		},
		{
			url:    "/value/gauge/test2",
			method: http.MethodGet,
			name:   "Get gauge test2",
			want: want{
				code:    http.StatusOK,
				content: "100.123",
			},
		},

		{
			url:    "/value/gauge/test3",
			method: http.MethodGet,
			name:   "Undefined metric",
			want: want{
				code:    http.StatusNotFound,
				content: "",
			},
		},
		{
			url:    "/value/gauge/test2",
			method: http.MethodPost,
			name:   "Try Post Method",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = test.method
			req.URL = server.URL + test.url

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			if len(test.want.content) > 0 {
				assert.Equal(t, string(resp.Body()), test.want.content)

			}

			assert.Equal(t, test.want.code, resp.StatusCode())
		})
	}
}

func TestGetMetricList(t *testing.T) {

	var store = MockMemoryStorage{
		Counters: map[string]storage.MetricCounterInfo{"first metric": {Name: "first metric", Value: 591}},
		Gauges:   map[string]storage.MetricGaugeInfo{"second metric": {Name: "second metric", Value: 100.123}},
	}

	r := MakeRouter(&store)
	server := httptest.NewServer(r)

	defer server.Close()

	t.Run("Get counter list", func(t *testing.T) {
		req := resty.New().R()
		req.Method = http.MethodGet
		req.URL = server.URL + "/"

		resp, err := req.Send()

		assert.NoError(t, err, "error making HTTP request")

		assert.Contains(t, string(resp.Body()), "first metric")
		assert.Contains(t, string(resp.Body()), "second metric")

		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})
}
