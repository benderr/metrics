package ipcheck_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benderr/metrics/pkg/ipcheck"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddleware_valid(t *testing.T) {

	m := ipcheck.Middleware("192.168.1.0/24")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		require.NoError(t, err)
	})

	handlerToTest := m(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)

	req.Header.Add("X-Real-IP", "192.168.1.134")

	r := httptest.NewRecorder()
	handlerToTest.ServeHTTP(r, req)
	resp := r.Result()
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, 200)
}

func TestMiddleware_noIP(t *testing.T) {

	m := ipcheck.Middleware("192.0.2.32/24")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		require.NoError(t, err)
	})

	handlerToTest := m(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)

	r := httptest.NewRecorder()
	handlerToTest.ServeHTTP(r, req)
	resp := r.Result()
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, 403)
}

func TestMiddleware_untrusted(t *testing.T) {

	m := ipcheck.Middleware("192.0.2.32/24")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		require.NoError(t, err)
	})

	handlerToTest := m(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	req.Header.Add("X-Real-IP", "192.0.3.25")

	r := httptest.NewRecorder()
	handlerToTest.ServeHTTP(r, req)
	resp := r.Result()
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, 403)
}

func TestMiddleware_panic(t *testing.T) {
	require.Panics(t, func() {
		ipcheck.Middleware("invalid subnet")
	})
}

func TestMiddleware_empty(t *testing.T) {

	m := ipcheck.Middleware("")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		require.NoError(t, err)
	})

	handlerToTest := m(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)

	r := httptest.NewRecorder()
	handlerToTest.ServeHTTP(r, req)
	resp := r.Result()
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, 200)
}
