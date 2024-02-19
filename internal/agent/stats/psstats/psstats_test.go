package psstats_test

import (
	"context"
	"testing"
	"time"

	"github.com/benderr/metrics/internal/agent/stats/psstats"
	"github.com/stretchr/testify/assert"
)

func TestPsStats_Collect(t *testing.T) {

	t.Run("should collected psstats", func(t *testing.T) {
		m := psstats.New(time.Millisecond * 1)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ch := m.Collect(ctx)
		val := <-ch

		assert.True(t, len(val) >= 2, "should contain min 2 object")
	})
}
