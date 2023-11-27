package sign

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

// Миддлвар для проверки подписи получаемого запроса
func (h *signValidator) CheckSign(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.secret == "" {
			next.ServeHTTP(w, r)
			return
		}
		hash := r.Header.Get("HashSHA256")
		if hash != "" {
			sign, err := hex.DecodeString(hash)

			if err != nil {
				h.logger.Errorln("decode hash error", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			buf := &bytes.Buffer{}

			teeReader := io.TeeReader(r.Body, buf)

			content, err := io.ReadAll(teeReader)

			if err != nil {
				h.logger.Errorln("can't read body", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			hinst := hmac.New(sha256.New, []byte(h.secret))
			hinst.Write(content)
			signFromBody := hinst.Sum(nil)

			if !hmac.Equal(sign, signFromBody) {
				h.logger.Infow("invalid sign")
				http.Error(w, "invalid sign", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(buf)
		}

		next.ServeHTTP(w, r)
	})
}
