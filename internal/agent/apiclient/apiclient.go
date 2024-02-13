package apiclient

import (
	"errors"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/benderr/metrics/internal/agent/logger"
	"github.com/benderr/metrics/pkg/sign"
)

// Расширяем апи рести под наши бизнес требования
type Client struct {
	*resty.Client
	secret string
	logger logger.Logger
}

func New(server string, secret string, logger logger.Logger) *Client {
	client := resty.
		New().
		SetBaseURL(server)

	return &Client{
		Client: client,
		secret: secret,
		logger: logger,
	}
}

const (
	ATTEMPT_1 int = 1
	ATTEMPT_2 int = 2
	ATTEMPT_3 int = 3
)

const (
	WAIT_1 int = 1
	WAIT_2 int = 3
	WAIT_3 int = 5
)

const MAX_RETRIES int64 = 5

// Кастомный конфиг для ретраев
func (a *Client) SetCustomRetries(count int) *Client {
	a.SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(time.Duration(MAX_RETRIES) * time.Second).
		SetRetryCount(count).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			wait := 0
			switch resp.Request.Attempt {
			case ATTEMPT_1:
				wait = WAIT_1
			case ATTEMPT_2:
				wait = WAIT_2
			case ATTEMPT_3:
				wait = WAIT_3
			}
			if wait > 0 {
				return time.Duration(wait) * time.Second, nil
			} else {
				return 0, errors.New("quota exceeded")
			}
		})
	return a
}

// Мидлвар для добавления подписанного ключом запроса
func (a *Client) SetSignedHeader() *Client {
	if a.secret != "" {
		a.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			if r.Header.Get("HashSHA256") == "" {
				if body, ok := a.getRequestBody(r); ok {
					signhex := sign.New(a.secret, body)
					a.logger.Infoln("generated sign", signhex)
					r.SetHeader("HashSHA256", signhex)
				}
			}
			return nil
		})
	}

	return a
}

func (a *Client) getRequestBody(r *resty.Request) ([]byte, bool) {
	if r.Body == nil {
		return []byte{}, false
	}
	if b, ok := r.Body.(string); ok {
		return []byte(b), true
	} else if b, ok := r.Body.([]byte); ok {
		return b, true
	}
	return []byte{}, false
}
