package psstats

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/benderr/metrics/internal/agent/stats"
)

type PsStats struct {
	interval time.Duration
}

// New return object for collect gopsutil metrics
func New(t time.Duration) *PsStats {
	return &PsStats{
		interval: t,
	}
}

func (m *PsStats) Collect(ctx context.Context) <-chan []stats.Item {
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

func (m *PsStats) getStats() []stats.Item {
	res := make([]stats.Item, 0)
	if v, err := mem.VirtualMemory(); err == nil {
		res = append(res, stats.Item{Type: "gauge", Name: "TotalMemory", Value: float64(v.Total)})
		res = append(res, stats.Item{Type: "gauge", Name: "FreeMemory", Value: float64(v.Free)})
	}

	if info, err := cpu.Percent(0, false); err == nil && len(info) > 0 {
		res = append(res, stats.Item{Type: "gauge", Name: "CPUutilization1", Value: info[0]})
	}

	return res
}
