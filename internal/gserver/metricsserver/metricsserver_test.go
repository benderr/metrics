package metricsserver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/benderr/metrics/internal/gserver/metricsserver"
	"github.com/benderr/metrics/internal/server/repository"
	"github.com/benderr/metrics/internal/server/repository/mockrepo"
	"github.com/benderr/metrics/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type mockLog struct{}

func (m *mockLog) Infoln(args ...interface{}) {
	fmt.Println(args...)
}

func (m *mockLog) Errorln(args ...interface{}) {
	fmt.Println(args...)
}

func (m *mockLog) Infow(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg)
}

func TestMetricServer(t *testing.T) {
	t.Run("should success update metric", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repoMock := mockrepo.NewMockMetricRepository(ctrl)

		s := metricsserver.New(repoMock, &mockLog{}, "")

		ctx := context.Background()
		var delta int64 = 10
		var delta2 int64 = 20
		repoMock.
			EXPECT().
			Update(ctx, repository.Metrics{ID: "test", MType: "counter", Delta: &delta}).
			Return(&repository.Metrics{ID: "test", MType: "counter", Delta: &delta2}, nil)

		res, err := s.UpdateMetric(ctx, &proto.UpdateMetricRequest{
			Metric: &proto.MetricItem{
				ID:    "test",
				MType: proto.MetricType_COUNTER,
				Delta: delta,
			},
		})

		require.NoError(t, err)
		assert.EqualValues(t, &proto.MetricItem{
			ID:    "test",
			MType: proto.MetricType_COUNTER,
			Delta: delta2,
		}, res.Metric)
	})

	t.Run("should success get metric", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repoMock := mockrepo.NewMockMetricRepository(ctrl)

		s := metricsserver.New(repoMock, &mockLog{}, "")

		ctx := context.Background()
		var value float64 = 10
		metricId := "test"
		repoMock.
			EXPECT().
			Get(ctx, metricId).
			Return(&repository.Metrics{ID: metricId, MType: "gauge", Value: &value}, nil)

		res, err := s.GetMetric(ctx, &proto.GetMetricRequest{
			ID: metricId,
		})

		require.NoError(t, err)
		assert.EqualValues(t, &proto.MetricItem{
			ID:    metricId,
			MType: proto.MetricType_GAUGE,
			Value: value,
		}, res.Metric)
	})

	t.Run("should not found for get metric", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repoMock := mockrepo.NewMockMetricRepository(ctrl)

		s := metricsserver.New(repoMock, &mockLog{}, "")

		ctx := context.Background()
		metricId := "test"
		repoMock.
			EXPECT().
			Get(ctx, metricId).
			Return(nil, nil)

		_, err := s.GetMetric(ctx, &proto.GetMetricRequest{
			ID: metricId,
		})

		assert.Error(t, err)

	})

	t.Run("should success bulk update metrics", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repoMock := mockrepo.NewMockMetricRepository(ctrl)

		s := metricsserver.New(repoMock, &mockLog{}, "")

		ctx := context.Background()
		var delta int64 = 10
		var value float64 = 10
		repoMock.
			EXPECT().
			BulkUpdate(ctx, []repository.Metrics{
				{ID: "test1", MType: "counter", Delta: &delta},
				{ID: "test2", MType: "gauge", Value: &value},
			}).
			Return(nil)

		_, err := s.BulkUpdateMetric(ctx, &proto.BulkUpdateMetricRequest{
			Metrics: []*proto.MetricItem{
				{
					ID:    "test1",
					MType: proto.MetricType_COUNTER,
					Delta: delta,
				},
				{
					ID:    "test2",
					MType: proto.MetricType_GAUGE,
					Value: value,
				},
			},
		})

		require.NoError(t, err)

	})

}
