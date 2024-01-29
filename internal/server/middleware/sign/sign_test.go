package sign_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benderr/metrics/internal/server/logger"
	"github.com/benderr/metrics/internal/server/middleware/sign"
	signer "github.com/benderr/metrics/internal/sign"
	"github.com/go-chi/chi"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

type testModel struct {
	ID   int
	Name string
}

func TestCheckSignWithoutSecret(t *testing.T) {

	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	r.Post("/check", checkHandler)

	req := resty.New().SetBaseURL(server.URL).R().SetHeader("Content-Type", "application/json")

	t.Run("should valid json in response", func(t *testing.T) {
		resp, err := req.SetBody(&testModel{ID: 1, Name: "Value string"}).Post("/check")
		assert.NoError(t, err, "error making HTTP request")
		assert.Equal(t, 200, resp.StatusCode())
		assert.JSONEq(t, `{"ID":1, "Name":"Value string"}`, string(resp.Body()))
	})
}

func TestCheckSign(t *testing.T) {

	logger, sync := logger.New()
	defer sync()

	r := chi.NewRouter()

	secret := "123"

	mwsign := sign.New(secret, logger)

	r.Use(mwsign.CheckSign)

	server := httptest.NewServer(r)
	defer server.Close()

	r.Post("/check", checkHandler)

	req := resty.New().SetBaseURL(server.URL).R().SetHeader("Content-Type", "application/json")

	t.Run("should valid json in response", func(t *testing.T) {

		m := &testModel{ID: 1, Name: "Value string"}
		mBytes, _ := json.Marshal(m)
		signhex := signer.New(secret, mBytes)

		resp, err := req.
			SetBody(mBytes).
			SetHeader("HashSHA256", signhex).
			Post("/check")

		assert.NoError(t, err, "error making HTTP request")
		assert.Equal(t, 200, resp.StatusCode())
		assert.JSONEq(t, `{"ID":1, "Name":"Value string"}`, string(resp.Body()))
	})

	t.Run("should invalid secret", func(t *testing.T) {

		m := &testModel{ID: 1, Name: "Value string"}
		mBytes, _ := json.Marshal(m)
		signhex := signer.New("invalid secret", mBytes)

		resp, err := req.
			SetBody(mBytes).
			SetHeader("HashSHA256", signhex).
			Post("/check")

		assert.NoError(t, err, "error making HTTP request")
		assert.Equal(t, 400, resp.StatusCode())
	})

	t.Run("should invalid payload", func(t *testing.T) {

		m := &testModel{ID: 1, Name: "Value string"}
		mBytes, _ := json.Marshal(m)
		signhex := signer.New(secret, mBytes)

		resp, err := req.
			SetBody(&testModel{ID: 2, Name: "Value string"}).
			SetHeader("HashSHA256", signhex).
			Post("/check")

		assert.NoError(t, err, "error making HTTP request")
		assert.Equal(t, 400, resp.StatusCode())
	})
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	var model testModel
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &model); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(&model)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
