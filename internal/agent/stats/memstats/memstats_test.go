package memstats_test

import (
	"context"
	"testing"
	"time"

	"github.com/benderr/metrics/internal/agent/stats"
	"github.com/benderr/metrics/internal/agent/stats/memstats"
	"github.com/stretchr/testify/assert"
)

func TestMemStats_Collect(t *testing.T) {
	t.Run("should collected memstats", func(t *testing.T) {
		m := memstats.New(time.Millisecond * 1)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ch := m.Collect(ctx)
		val := <-ch

		assert.True(t, len(val) == 29, "should be 29")
		assert.Contains(t, val, stats.Item{Type: "counter", Name: "PollCount", Delta: 1}, "should contain PollCount")
	})
}
