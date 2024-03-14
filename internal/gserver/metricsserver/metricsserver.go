package metricsserver

import (
	"context"

	"github.com/benderr/metrics/internal/server/repository"
	"github.com/benderr/metrics/pkg/logger"
	pb "github.com/benderr/metrics/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type metricsServer struct {
	pb.UnimplementedMetricsServer
	metricRepo repository.MetricRepository
	logger     logger.Logger
	secret     string
}

func New(repo repository.MetricRepository, logger logger.Logger, secret string) *metricsServer {
	return &metricsServer{
		metricRepo: repo,
		logger:     logger,
		secret:     secret,
	}
}

func (m *metricsServer) UpdateMetric(ctx context.Context, in *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	if in.Metric == nil {
		return nil, status.Error(codes.InvalidArgument, "metric is nil")
	}
	mtr := convertToMetric(in.Metric)
	res, err := m.metricRepo.Update(ctx, mtr)

	if err != nil {
		m.logger.Infow("update metric", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	m.logger.Infoln("SUCCESS 1")
	return &pb.UpdateMetricResponse{
		Metric: convertToResponseMetric(res),
	}, nil
}

func (m *metricsServer) BulkUpdateMetric(ctx context.Context, in *pb.BulkUpdateMetricRequest) (*pb.BulkUpdateMetricResponse, error) {
	metrics := make([]repository.Metrics, len(in.Metrics))

	for i, mtr := range in.Metrics {
		metrics[i] = convertToMetric(mtr)
	}

	err := m.metricRepo.BulkUpdate(ctx, metrics)
	if err != nil {
		m.logger.Infow("bulk update metric", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	m.logger.Infoln("SUCCESS 2")
	return nil, nil
}

func (m *metricsServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {

	res, err := m.metricRepo.Get(ctx, in.ID)

	if err != nil {
		m.logger.Infow("get metric", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	m.logger.Infoln("SUCCESS 3")

	if res == nil {
		m.logger.Infow("not found metric")
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &pb.GetMetricResponse{
		Metric: convertToResponseMetric(res),
	}, nil
}

func (m *metricsServer) ListMetric(ctx context.Context, in *pb.ListMetricRequest) (*pb.ListMetricResponse, error) {

	res, err := m.metricRepo.GetList(ctx)

	if err != nil {
		m.logger.Infow("get metric", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	m.logger.Infoln("SUCCESS 4")

	metrics := make([]*pb.MetricItem, len(res))

	for i, mtr := range res {
		metrics[i] = convertToResponseMetric(&mtr)
	}

	return &pb.ListMetricResponse{
		Metrics: metrics,
	}, nil
}

func convertToMetric(m *pb.MetricItem) repository.Metrics {
	res := repository.Metrics{
		ID: m.ID,
	}
	switch m.MType {
	case pb.MetricType_COUNTER:
		res.MType = "counter"
		res.Delta = &m.Delta
	case pb.MetricType_GAUGE:
		res.MType = "gauge"
		res.Value = &m.Value
	}
	return res
}

func convertToResponseMetric(m *repository.Metrics) *pb.MetricItem {
	res := &pb.MetricItem{
		ID: m.ID,
	}
	switch m.MType {
	case "counter":
		res.MType = pb.MetricType_COUNTER
		res.Delta = *m.Delta
	case "gauge":
		res.MType = pb.MetricType_GAUGE
		res.Value = *m.Value
	}
	return res
}
