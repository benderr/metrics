package allstats_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/benderr/metrics/internal/agent/stats"
	"github.com/benderr/metrics/internal/agent/stats/allstats"
)

func TestJoin(t *testing.T) {

	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	c1 := &testCollector{items: []stats.Item{
		{Name: "test1", Type: "counter", Delta: 1},
		{Name: "test2", Type: "gauge", Value: 1.1},
	}}

	c2 := &testCollector{items: []stats.Item{
		{Name: "test3", Type: "counter", Delta: 2},
		{Name: "test4", Type: "gauge", Value: 2.1},
	}}

	ch := allstats.Join(ctxWithCancel, c1, c2)

	res := make([]stats.Item, 0)

	for i := range ch {
		res = append(res, i...)
	}

	assert.Equal(t, len(res), 4)
}

type testCollector struct {
	items []stats.Item
}

func (t *testCollector) Collect(ctx context.Context) <-chan []stats.Item {
	ch := make(chan []stats.Item)

	go func() {
		for _, item := range t.items {
			ch <- []stats.Item{item}
		}
		defer close(ch)
	}()

	return ch
}
