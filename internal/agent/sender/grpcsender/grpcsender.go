package grpcsender

import (
	"context"
	"errors"
	"log"

	"github.com/benderr/metrics/internal/agent/report"
	"github.com/benderr/metrics/pkg/logger"
	pb "github.com/benderr/metrics/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCSender struct {
	bulk   bool
	log    logger.Logger
	client pb.MetricsClient
}

func New(host string, logger logger.Logger, bulk bool) (*GRPCSender, func()) {
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := pb.NewMetricsClient(conn)
	cancel := func() {
		conn.Close()
	}
	return &GRPCSender{
		log:    logger,
		client: c,
		bulk:   bulk,
	}, cancel
}

func (g *GRPCSender) Send(ctx context.Context, metrics []report.MetricItem) error {
	if g.bulk {
		list := make([]*pb.MetricItem, len(metrics))
		for i, m := range metrics {
			mtr := convertToMetric(m)
			list[i] = &mtr
		}

		_, err := g.client.BulkUpdateMetric(ctx, &pb.BulkUpdateMetricRequest{
			Metrics: list,
		})
		g.log.Errorln("grpc bulk mode error: ", err)
		return err

	} else {
		errs := make([]error, 0)
		for _, m := range metrics {
			mtr := convertToMetric(m)
			_, err := g.client.UpdateMetric(ctx, &pb.UpdateMetricRequest{
				Metric: &mtr,
			})
			if err != nil {
				g.log.Errorln("grpc single mode error: ", err)
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return errors.Join(errs...)
		}
		return nil
	}
}

func convertToMetric(m report.MetricItem) pb.MetricItem {
	var delta int64
	var value float64
	if m.Delta != nil {
		delta = *m.Delta
	}

	if m.Value != nil {
		value = *m.Value
	}

	return pb.MetricItem{
		ID:    m.ID,
		MType: getMetricType(m.MType),
		Delta: delta,
		Value: value,
	}
}

func getMetricType(mType string) pb.MetricType {
	switch mType {
	case "counter":
		return pb.MetricType_COUNTER
	case "gauge":
		return pb.MetricType_GAUGE
	}
	return pb.MetricType_UNSPECIFIED
}
