package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func New(secret string, body []byte) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	signedBody := h.Sum(nil)
	return hex.EncodeToString(signedBody)
}
