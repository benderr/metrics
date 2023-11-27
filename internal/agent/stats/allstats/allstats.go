package allstats

import (
	"context"

	"github.com/benderr/metrics/internal/agent/stats"
)

type Collector interface {
	Collect(ctx context.Context) <-chan []stats.Item
}

func Join(ctx context.Context, collectors ...Collector) <-chan []stats.Item {
	outCh := make(chan []stats.Item)
	go func() {
		defer close(outCh)
		for _, c := range collectors {
			inChan := c.Collect(ctx)
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case v := <-inChan:
						outCh <- v
					}
				}
			}()
		}
	}()

	return outCh
}
