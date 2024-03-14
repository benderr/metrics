package gziper_test

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benderr/metrics/pkg/gziper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testModel struct {
	ID   string
	Name string
}

func TestGzipCompressor_TransformReader(t *testing.T) {

	m := gziper.New(1, "application/json")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		require.NoError(t, err)

		model := &testModel{}
		err = json.Unmarshal(buf.Bytes(), model)

		require.NoError(t, err)

		assert.Equal(t, model.ID, "1")
		assert.Equal(t, model.Name, "Test")
	})

	handlerToTest := m.TransformReader(nextHandler)

	v, err := json.Marshal(&testModel{ID: "1", Name: "Test"})

	require.NoError(t, err)

	vCompressed, err := compress(v)

	require.NoError(t, err)

	req := httptest.NewRequest("GET", "http://testing", bytes.NewBuffer(vCompressed)) //Content-Encoding

	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Content-Type", "application/json")

	r := httptest.NewRecorder()
	handlerToTest.ServeHTTP(r, req)
	resp := r.Result()
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, 200)

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
