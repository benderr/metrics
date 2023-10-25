package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benderr/metrics/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

type MockMemoryStorage struct {
	Metrics map[string]storage.Metrics
}

func (m *MockMemoryStorage) Update(mtr storage.Metrics) (*storage.Metrics, error) {
	if metric, ok := m.Metrics[mtr.ID]; ok {
		updatedMetric := storage.Metrics{
			ID:    metric.ID,
			MType: metric.MType,
		}
		switch metric.MType {
		case "counter":
			newVal := *metric.Delta + *mtr.Delta
			updatedMetric.Delta = &newVal
		case "gauge":
			updatedMetric.Value = mtr.Value
		}
		m.Metrics[mtr.ID] = updatedMetric
		return &updatedMetric, nil
	} else {
		m.Metrics[mtr.ID] = mtr
		res := m.Metrics[mtr.ID]
		return &res, nil
	}

}

func (m *MockMemoryStorage) GetList() ([]storage.Metrics, error) {
	res := []storage.Metrics{}
	for _, item := range m.Metrics {
		res = append(res, item)
	}
	return res, nil
}

func (m *MockMemoryStorage) Get(name string) (*storage.Metrics, error) {
	if res, ok := m.Metrics[name]; ok {
		return &storage.Metrics{
			ID:    res.ID,
			Value: res.Value,
			Delta: res.Delta,
			MType: res.MType,
		}, nil
	}
	return nil, nil
}

func TestUpdateMetricByUrlHandler(t *testing.T) {

	var store = MockMemoryStorage{
		Metrics: make(map[string]storage.Metrics),
	}

	h := NewHandlers(&store)
	r := chi.NewRouter()
	h.AddHandlers(r)
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
			url:    "/update/gauge/test3/2.0",
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

func TestGetMetricByUrlHandler(t *testing.T) {

	var delta int64 = 1
	var val1 float64 = 100.1200
	var val2 float64 = 806132.0

	var store = MockMemoryStorage{
		Metrics: map[string]storage.Metrics{
			"test":   {ID: "test", Delta: &delta, MType: "counter"},
			"test2":  {ID: "test2", Value: &val1, MType: "gauge"},
			"test22": {ID: "test22", Value: &val2, MType: "gauge"},
		},
	}

	h := NewHandlers(&store)
	r := chi.NewRouter()
	h.AddHandlers(r)
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
				content: "100.12",
			},
		},
		{
			url:    "/value/gauge/test22",
			method: http.MethodGet,
			name:   "Get gauge test22",
			want: want{
				code:    http.StatusOK,
				content: "806132.",
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

	var delta int64 = 591
	var val1 float64 = 100.1200

	var store = MockMemoryStorage{
		Metrics: map[string]storage.Metrics{
			"test":   {ID: "first metric", Delta: &delta, MType: "counter"},
			"test22": {ID: "second metric", Value: &val1, MType: "gauge"},
		},
	}

	h := NewHandlers(&store)
	r := chi.NewRouter()
	h.AddHandlers(r)
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

func TestGetMetricHandler(t *testing.T) {

	var delta int64 = 1
	var val1 float64 = 100.1200
	var val2 float64 = 806132.0

	var store = MockMemoryStorage{
		Metrics: map[string]storage.Metrics{
			"test":   {ID: "test", Delta: &delta, MType: "counter"},
			"test2":  {ID: "test2", Value: &val1, MType: "gauge"},
			"test22": {ID: "test22", Value: &val2, MType: "gauge"},
		},
	}

	h := NewHandlers(&store)
	r := chi.NewRouter()
	h.AddHandlers(r)
	server := httptest.NewServer(r)

	defer server.Close()

	type want struct {
		code    int
		content string
	}
	tests := []struct {
		name string
		body *storage.Metrics
		url  string
		want want
	}{
		{
			url:  "/value/",
			name: "Get counter test",
			body: &storage.Metrics{
				ID:    "test",
				MType: "gauge",
			},
			want: want{
				code:    http.StatusOK,
				content: `{"delta":1, "id":"test", "type":"counter"}`,
			},
		},
		{
			url:  "/value/",
			name: "Get gauge test2",
			body: &storage.Metrics{
				ID:    "test2",
				MType: "gauge",
			},
			want: want{
				code:    http.StatusOK,
				content: `{"value":100.12, "id":"test2", "type":"gauge"}`,
			},
		},
	}

	req := resty.New().SetBaseURL(server.URL).R().SetHeader("Content-Type", "application/json")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := req.
				SetBody(test.body).
				Post(test.url)

			assert.NoError(t, err, "error making HTTP request")

			if len(test.want.content) > 0 {
				assert.JSONEq(t, string(resp.Body()), test.want.content)
			}

			assert.Equal(t, test.want.code, resp.StatusCode())
		})
	}
}

func TestUpdateMetricHandler(t *testing.T) {
	var delta int64 = 1
	var val1 float64 = 100.1200
	var val2 float64 = 806132.0

	var store = MockMemoryStorage{
		Metrics: map[string]storage.Metrics{
			"test":   {ID: "test", Delta: &delta, MType: "counter"},
			"test2":  {ID: "test2", Value: &val1, MType: "gauge"},
			"test22": {ID: "test22", Value: &val2, MType: "gauge"},
		},
	}

	h := NewHandlers(&store)
	r := chi.NewRouter()
	h.AddHandlers(r)
	server := httptest.NewServer(r)

	defer server.Close()

	var resDelta int64 = 2
	var resValue float64 = 102.1200

	type want struct {
		code    int
		content string
	}
	tests := []struct {
		name string
		url  string
		body *storage.Metrics
		want want
	}{
		{
			url: "/update",
			body: &storage.Metrics{
				ID:    "test",
				MType: "counter",
				Delta: &resDelta,
			},
			name: "Add new counter metric",
			want: want{
				code:    http.StatusOK,
				content: `{"delta":3, "id":"test", "type":"counter"}`,
			},
		},
		{
			url: "/update",
			body: &storage.Metrics{
				ID:    "test2",
				MType: "gauge",
				Value: &resValue,
			},
			name: "Add new gauge metric",
			want: want{
				code:    http.StatusOK,
				content: `{"value":102.12, "id":"test2", "type":"gauge"}`,
			},
		},
	}

	req := resty.New().SetBaseURL(server.URL).R().SetHeader("Content-Type", "application/json")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := req.
				SetBody(test.body).
				Post(test.url)

			assert.NoError(t, err, "error making HTTP request")

			if len(test.want.content) > 0 {
				assert.JSONEq(t, string(resp.Body()), test.want.content)
			}

			assert.Equal(t, test.want.code, resp.StatusCode())
		})
	}
}
