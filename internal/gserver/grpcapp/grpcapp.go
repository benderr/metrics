package grpcapp

import (
	"context"
	"net"

	"github.com/benderr/metrics/internal/gserver/interceptors/ilog"
	"github.com/benderr/metrics/internal/gserver/metricsserver"
	"github.com/benderr/metrics/internal/server/config"
	"github.com/benderr/metrics/internal/server/repository/storage"
	"github.com/benderr/metrics/pkg/logger"
	pb "github.com/benderr/metrics/proto"
	"google.golang.org/grpc"
)

type grpcApp struct {
	config *config.Config
	log    logger.Logger
}

func New(config *config.Config, log logger.Logger) *grpcApp {
	return &grpcApp{
		config: config,
		log:    log,
	}
}

// Run create and run grpc server
func (a *grpcApp) Run(ctx context.Context) error {

	repo, err := storage.New(ctx, a.config, a.log)
	if err != nil {
		return err
	}

	listen, err := net.Listen("tcp", string(a.config.Server))

	if err != nil {
		return err
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(ilog.New(a.log).Interceptor))

	srv := metricsserver.New(repo, a.log, a.config.SecretKey)

	pb.RegisterMetricsServer(s, srv)

	if err := s.Serve(listen); err != nil {
		return err
	}

	return nil
}
