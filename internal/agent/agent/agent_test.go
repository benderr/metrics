package agent_test

import (
	"context"
	"testing"

	"github.com/benderr/metrics/internal/agent/agent"
	"github.com/benderr/metrics/internal/agent/agent/mocks"
	"github.com/benderr/metrics/internal/agent/report"
	mocksender "github.com/benderr/metrics/internal/agent/sender/mocks"
	"github.com/benderr/metrics/internal/agent/stats"
	"go.uber.org/mock/gomock"
)

func TestAgent_Run(t *testing.T) {

	t.Run("should update report", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reportMock := mocks.NewMockIReport(ctrl)
		senderMock := mocksender.NewMockMetricSender(ctrl)

		a := agent.New(senderMock, reportMock)

		ch := make(chan struct{})
		in := make(chan []stats.Item)
		defer close(in)
		defer close(ch)

		metricList := []stats.Item{{Name: "test", Value: 1}}
		ctx, cancel := context.WithCancel(context.Background()) //with cancel to trigger stop

		go func() {
			in <- metricList
			cancel()
		}()

		reportMock.EXPECT().Update(metricList)

		a.Run(ctx, in, ch)
	})

	t.Run("should send report", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reportMock := mocks.NewMockIReport(ctrl)
		senderMock := mocksender.NewMockMetricSender(ctrl)

		a := agent.New(senderMock, reportMock)

		ch := make(chan struct{})
		in := make(chan []stats.Item)
		defer close(ch)
		defer close(in)

		var delta float64 = 1
		sendList := []report.MetricItem{{ID: "test", Value: &delta}}

		ctx, cancel := context.WithCancel(context.Background()) //with cancel to trigger stop

		go func() {
			ch <- struct{}{}
			cancel()
		}()

		reportMock.EXPECT().GetList().Return(sendList)
		senderMock.EXPECT().Send(ctx, sendList).Return(nil)

		a.Run(ctx, in, ch)
	})
}
