package handlers_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/benderr/metrics/internal/server/handlers"
	"github.com/benderr/metrics/internal/server/repository"
	"github.com/benderr/metrics/pkg/gziper"
)

type MockMemoryStorage struct {
	Metrics map[string]repository.Metrics
}

type MockLogger struct{}

func (m *MockLogger) Infoln(args ...interface{}) {
	fmt.Println(args...)
}

func (m *MockLogger) Errorln(args ...interface{}) {
	fmt.Println(args...)
}

func (m *MockLogger) Infow(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg)
}

func (m *MockMemoryStorage) Update(ctx context.Context, mtr repository.Metrics) (*repository.Metrics, error) {
	if metric, ok := m.Metrics[mtr.ID]; ok {
		updatedMetric := repository.Metrics{
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

func (m *MockMemoryStorage) BulkUpdate(ctx context.Context, metrics []repository.Metrics) error {

	if len(metrics) == 0 {
		return nil
	}

	for _, v := range metrics {
		m.Update(ctx, v)
	}

	return nil
}

func (m *MockMemoryStorage) GetList(ctx context.Context) ([]repository.Metrics, error) {
	res := []repository.Metrics{}
	for _, item := range m.Metrics {
		res = append(res, item)
	}
	return res, nil
}

func (m *MockMemoryStorage) Get(ctx context.Context, name string) (*repository.Metrics, error) {
	if res, ok := m.Metrics[name]; ok {
		return &repository.Metrics{
			ID:    res.ID,
			Value: res.Value,
			Delta: res.Delta,
			MType: res.MType,
		}, nil
	}
	return nil, nil
}

func (m *MockMemoryStorage) PingContext(ctx context.Context) error {
	return nil
}

func TestUpdateMetricByUrlHandler(t *testing.T) {

	var store = MockMemoryStorage{
		Metrics: make(map[string]repository.Metrics),
	}

	h := handlers.New(&store, &MockLogger{}, "")
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
	val1 := 100.1200
	val2 := 806132.0

	var store = MockMemoryStorage{
		Metrics: map[string]repository.Metrics{
			"test":   {ID: "test", Delta: &delta, MType: "counter"},
			"test2":  {ID: "test2", Value: &val1, MType: "gauge"},
			"test22": {ID: "test22", Value: &val2, MType: "gauge"},
		},
	}

	h := handlers.New(&store, &MockLogger{}, "")
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
	val1 := 100.1200

	var store = MockMemoryStorage{
		Metrics: map[string]repository.Metrics{
			"test":   {ID: "first metric", Delta: &delta, MType: "counter"},
			"test22": {ID: "second metric", Value: &val1, MType: "gauge"},
		},
	}

	h := handlers.New(&store, &MockLogger{}, "")
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
	val1 := 100.1200

	var store = MockMemoryStorage{
		Metrics: map[string]repository.Metrics{
			"test":  {ID: "test", Delta: &delta, MType: "counter"},
			"test2": {ID: "test2", Value: &val1, MType: "gauge"},
		},
	}

	h := handlers.New(&store, &MockLogger{}, "")
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
		body *repository.Metrics
		url  string
		want want
	}{
		{
			url:  "/value/",
			name: "Get counter test",
			body: &repository.Metrics{
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
			body: &repository.Metrics{
				ID:    "test2",
				MType: "gauge",
			},
			want: want{
				code:    http.StatusOK,
				content: `{"value":100.12, "id":"test2", "type":"gauge"}`,
			},
		},
		{
			url:  "/value/",
			name: "Get not exist metric",
			body: &repository.Metrics{
				ID:    "test3",
				MType: "gauge",
			},
			want: want{
				code: http.StatusNotFound,
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
	val1 := 100.1200
	val2 := 806132.0

	var store = MockMemoryStorage{
		Metrics: map[string]repository.Metrics{
			"test":   {ID: "test", Delta: &delta, MType: "counter"},
			"test2":  {ID: "test2", Value: &val1, MType: "gauge"},
			"test22": {ID: "test22", Value: &val2, MType: "gauge"},
		},
	}

	h := handlers.New(&store, &MockLogger{}, "")
	r := chi.NewRouter()
	h.AddHandlers(r)
	server := httptest.NewServer(r)

	defer server.Close()

	var resDelta int64 = 2
	resValue := 102.1200

	type want struct {
		code    int
		content string
	}
	tests := []struct {
		name string
		url  string
		body *repository.Metrics
		want want
	}{
		{
			url: "/update",
			body: &repository.Metrics{
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
			body: &repository.Metrics{
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

func TestGetMetricAcceptGzipOutputHandler(t *testing.T) {

	val1 := 100.1200

	var store = MockMemoryStorage{
		Metrics: map[string]repository.Metrics{
			"test2": {ID: "test2", Value: &val1, MType: "gauge"},
		},
	}

	h := handlers.New(&store, &MockLogger{}, "")
	r := chi.NewRouter()
	g := gziper.New(1, "application/json", "text/html")
	r.Use(g.TransformWriter)
	h.AddHandlers(r)

	server := httptest.NewServer(r)

	defer server.Close()

	req := resty.
		New().
		SetBaseURL(server.URL).
		R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip")

	t.Run("Get gauge test2 with accept-encoding: gzip", func(t *testing.T) {
		resp, err := req.
			SetBody(&repository.Metrics{
				ID:    "test2",
				MType: "gauge",
			}).
			Post("/value/")

		assert.NoError(t, err, "error making HTTP request")

		assert.JSONEq(t, string(resp.Body()), `{"value":100.12, "id":"test2", "type":"gauge"}`)

		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})
}

func TestGetMetricAcceptGzipInputHandler(t *testing.T) {

	val1 := 100.1200

	var store = MockMemoryStorage{
		Metrics: map[string]repository.Metrics{
			"test2": {ID: "test2", Value: &val1, MType: "gauge"},
		},
	}

	h := handlers.New(&store, &MockLogger{}, "")
	r := chi.NewRouter()
	g := gziper.New(1, "application/json", "text/html")
	r.Use(g.TransformWriter)
	r.Use(g.TransformReader)

	h.AddHandlers(r)

	server := httptest.NewServer(r)

	defer server.Close()

	req := resty.
		New().
		SetBaseURL(server.URL).
		R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip")

	t.Run("Get gauge test2 with gzipped request body", func(t *testing.T) {

		body, err := compress([]byte(`{"value":100.12, "id":"test2", "type":"gauge"}`))
		require.NoError(t, err)

		resp, err := req.
			SetBody(body).
			Post("/value/")

		assert.NoError(t, err, "error making HTTP request")
		assert.JSONEq(t, string(resp.Body()), `{"value":100.12, "id":"test2", "type":"gauge"}`)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})
}

func TestBulkUpdateHandler(t *testing.T) {
	var store = MockMemoryStorage{
		Metrics: make(map[string]repository.Metrics),
	}

	h := handlers.New(&store, &MockLogger{}, "")
	r := chi.NewRouter()
	g := gziper.New(1, "application/json", "text/html")
	r.Use(g.TransformWriter)
	r.Use(g.TransformReader)
	h.AddHandlers(r)

	server := httptest.NewServer(r)

	defer server.Close()

	type want struct {
		code int
	}

	var delta int64 = 1
	value := 100.1200

	tests := []struct {
		name string
		url  string
		body []repository.Metrics
		want want
	}{
		{
			url: "/updates/",
			body: []repository.Metrics{
				{ID: "test", MType: "counter", Delta: &delta},
				{ID: "test2", MType: "gauge", Value: &value},
			},
			name: "Add new metric list",
			want: want{
				code: http.StatusOK,
			},
		},
	}

	req := resty.New().SetBaseURL(server.URL).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			metrics, _ := json.Marshal(test.body)

			body, err := compress(metrics)

			require.NoError(t, err)

			resp, err := req.
				SetBody(body).
				Post(test.url)

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, test.want.code, resp.StatusCode())
		})
	}
}

func TestParseCounter(t *testing.T) {
	t.Run("should parse counter success", func(t *testing.T) {
		m, err := handlers.ParseCounter("counter", "test", "10")
		assert.NoError(t, err)
		assert.Equal(t, *m.Delta, int64(10))
		assert.Equal(t, m.MType, "counter")
		assert.Equal(t, m.ID, "test")
	})

	t.Run("should parse counter error delta", func(t *testing.T) {
		_, err := handlers.ParseCounter("counter", "test", "10b")
		assert.Error(t, err, "invalid delta")
	})

	t.Run("should parse counter error", func(t *testing.T) {
		_, err := handlers.ParseCounter("counter1", "test", "10b")
		assert.Error(t, err, "invalid metric type")
	})
}

func TestParseGauge(t *testing.T) {
	t.Run("should parse gauge success", func(t *testing.T) {
		m, err := handlers.ParseGauge("gauge", "test", "10.1")
		assert.NoError(t, err)
		assert.Equal(t, *m.Value, float64(10.1))
		assert.Equal(t, m.MType, "gauge")
		assert.Equal(t, m.ID, "test")
	})

	t.Run("should parse gauge error value", func(t *testing.T) {
		_, err := handlers.ParseCounter("counter", "test", "10b")
		assert.Error(t, err, "invalid value")
	})

	t.Run("should parse gauge error", func(t *testing.T) {
		_, err := handlers.ParseCounter("gauge1", "test", "10b")
		assert.Error(t, err, "invalid metric type")
	})
}

func BenchmarkGetMetricHandler(b *testing.B) {
	var delta int64 = 1
	val1 := 100.1200

	var store = MockMemoryStorage{
		Metrics: map[string]repository.Metrics{
			"test":  {ID: "test", Delta: &delta, MType: "counter"},
			"test2": {ID: "test2", Value: &val1, MType: "gauge"},
		},
	}
	h := handlers.New(&store, &MockLogger{}, "")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		v, _ := json.Marshal(&repository.Metrics{
			ID:    "test",
			MType: "gauge",
		})
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/", bytes.NewBuffer(v))
		rw := httptest.NewRecorder()
		b.StartTimer()
		h.GetMetricHandler(rw, req)
	}
}

func compress(s []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	zipped := gzip.NewWriter(buf)
	_, err := zipped.Write(s)
	if err != nil {
		return nil, err
	}
	zipped.Close()
	return buf.Bytes(), nil
}
