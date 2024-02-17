package ticker_test

import (
	"context"
	"testing"
	"time"

	"github.com/benderr/metrics/internal/agent/ticker"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5)
	defer cancel()
	ch := ticker.New(ctx, time.Millisecond*2)

	select {
	case <-ctx.Done():
		assert.Fail(t, "ticker don't ticked")
	case <-ch:
		t.Log("ticked 1")
	}
}
