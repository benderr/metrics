package sign

import (
	"github.com/benderr/metrics/internal/server/logger"
)

func New(secret string, logger logger.Logger) *signValidator {
	return &signValidator{
		secret: secret,
		logger: logger,
	}
}

type signValidator struct {
	secret string
	logger logger.Logger
}
