package ilog

import (
	"context"
	"fmt"

	"github.com/benderr/metrics/pkg/logger"
	"google.golang.org/grpc"
)

type logInterceptor struct {
	logger logger.Logger
}

func New(l logger.Logger) *logInterceptor {
	return &logInterceptor{
		logger: l,
	}
}

func (l *logInterceptor) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	res, err := handler(ctx, req)

	l.logger.Infoln(
		"method", info.FullMethod,
		"size", res,
	)

	if err != nil {
		l.logger.Errorln(fmt.Sprintf("method: %v, error: %v", info.FullMethod, err))
	}

	return res, err
}
