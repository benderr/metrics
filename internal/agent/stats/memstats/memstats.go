package memstats

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	"github.com/benderr/metrics/internal/agent/stats"
)

type MemStats struct {
	interval time.Duration
}

// New return memStats for collect runtime metrics
func New(t time.Duration) *MemStats {
	return &MemStats{
		interval: t,
	}
}

func (m *MemStats) Collect(ctx context.Context) <-chan []stats.Item {
	outCh := make(chan []stats.Item)
	go func() {
		defer close(outCh)
		pollTicker := time.NewTicker(m.interval)
		defer pollTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-pollTicker.C:
				outCh <- m.getStats()
			}
		}
	}()

	return outCh
}

func (m *MemStats) getStats() []stats.Item {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	res := make([]stats.Item, 0)

	res = append(res, stats.Item{Type: "gauge", Name: "Alloc", Value: float64(rtm.Alloc)})
	res = append(res, stats.Item{Type: "gauge", Name: "BuckHashSys", Value: float64(rtm.BuckHashSys)})
	res = append(res, stats.Item{Type: "gauge", Name: "Frees", Value: float64(rtm.Frees)})
	res = append(res, stats.Item{Type: "gauge", Name: "GCCPUFraction", Value: float64(rtm.GCCPUFraction)})
	res = append(res, stats.Item{Type: "gauge", Name: "GCSys", Value: float64(rtm.GCSys)})
	res = append(res, stats.Item{Type: "gauge", Name: "HeapAlloc", Value: float64(rtm.HeapAlloc)})
	res = append(res, stats.Item{Type: "gauge", Name: "HeapIdle", Value: float64(rtm.HeapIdle)})
	res = append(res, stats.Item{Type: "gauge", Name: "HeapInuse", Value: float64(rtm.HeapInuse)})
	res = append(res, stats.Item{Type: "gauge", Name: "HeapObjects", Value: float64(rtm.HeapObjects)})
	res = append(res, stats.Item{Type: "gauge", Name: "HeapReleased", Value: float64(rtm.HeapReleased)})
	res = append(res, stats.Item{Type: "gauge", Name: "HeapSys", Value: float64(rtm.HeapSys)})
	res = append(res, stats.Item{Type: "gauge", Name: "LastGC", Value: float64(rtm.LastGC)})
	res = append(res, stats.Item{Type: "gauge", Name: "Lookups", Value: float64(rtm.Lookups)})
	res = append(res, stats.Item{Type: "gauge", Name: "MCacheInuse", Value: float64(rtm.MCacheInuse)})
	res = append(res, stats.Item{Type: "gauge", Name: "MCacheSys", Value: float64(rtm.MCacheSys)})
	res = append(res, stats.Item{Type: "gauge", Name: "MSpanInuse", Value: float64(rtm.MSpanInuse)})
	res = append(res, stats.Item{Type: "gauge", Name: "MSpanSys", Value: float64(rtm.MSpanSys)})
	res = append(res, stats.Item{Type: "gauge", Name: "Mallocs", Value: float64(rtm.Mallocs)})
	res = append(res, stats.Item{Type: "gauge", Name: "NextGC", Value: float64(rtm.NextGC)})
	res = append(res, stats.Item{Type: "gauge", Name: "NumForcedGC", Value: float64(rtm.NumForcedGC)})
	res = append(res, stats.Item{Type: "gauge", Name: "NumGC", Value: float64(rtm.NumGC)})
	res = append(res, stats.Item{Type: "gauge", Name: "OtherSys", Value: float64(rtm.OtherSys)})
	res = append(res, stats.Item{Type: "gauge", Name: "PauseTotalNs", Value: float64(rtm.PauseTotalNs)})
	res = append(res, stats.Item{Type: "gauge", Name: "StackInuse", Value: float64(rtm.StackInuse)})
	res = append(res, stats.Item{Type: "gauge", Name: "StackSys", Value: float64(rtm.StackSys)})
	res = append(res, stats.Item{Type: "gauge", Name: "Sys", Value: float64(rtm.Sys)})
	res = append(res, stats.Item{Type: "gauge", Name: "TotalAlloc", Value: float64(rtm.TotalAlloc)})
	res = append(res, stats.Item{Type: "gauge", Name: "RandomValue", Value: rand.Float64()})
	res = append(res, stats.Item{Type: "counter", Name: "PollCount", Delta: 1})
	return res
}
