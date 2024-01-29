package allstats

import (
	"context"
	"sync"

	"github.com/benderr/metrics/internal/agent/stats"
)

type ICollect interface {
	Collect(ctx context.Context) <-chan []stats.Item
}

func Join(ctx context.Context, collectors ...ICollect) <-chan []stats.Item {
	outCh := make(chan []stats.Item)
	go func() {
		wg := sync.WaitGroup{}

		for _, c := range collectors {
			wg.Add(1)
			inChan := c.Collect(ctx)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					case v, ok := <-inChan:
						if !ok {
							return
						}
						outCh <- v
					}
				}
			}()
		}

		go func() {
			wg.Wait() //закрываем канал если вышли из всех горутин
			close(outCh)
		}()

	}()

	return outCh
}
