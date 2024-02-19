package ticker

import (
	"context"
	"time"
)

// New return channel which send message with interval
func New(ctx context.Context, interval time.Duration) <-chan struct{} {
	t := time.NewTicker(interval)
	ch := make(chan struct{})

	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				t.Stop()
			case <-t.C:
				ch <- struct{}{}
			}
		}
	}()
	return ch
}
