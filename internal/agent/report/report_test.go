package report_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/benderr/metrics/internal/agent/report"
	"github.com/benderr/metrics/internal/agent/stats"
)

func TestReport(t *testing.T) {
	t.Run("Test report read and write", func(t *testing.T) {

		r := report.New()

		sl := make([]stats.Item, 0)
		sl = append(sl, stats.Item{Name: "test", Type: "gauge", Value: 100.12})
		sl = append(sl, stats.Item{Name: "test2", Type: "counter", Delta: 1})

		r.Update(sl)

		res := r.GetList()

		for _, m := range res {
			switch m.MType {
			case "gauge":
				assert.Equal(t, *m.Value, 100.12)
			case "counter":
				var delta int64 = 1
				assert.Equal(t, *m.Delta, delta)
			}
		}

		assert.Equal(t, len(res), 2)
	})
}
