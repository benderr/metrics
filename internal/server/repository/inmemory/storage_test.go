package inmemory_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/benderr/metrics/internal/server/repository"
	"github.com/benderr/metrics/internal/server/repository/inmemory"
)

func BenchmarkGet(b *testing.B) {

	s1 := inmemory.New()
	s2 := inmemory.NewFast()

	ctx := context.Background()

	// write mock items to storage
	var delta int64 = 1
	val1 := 100.1200
	for i := 0; i < 1000; i++ {
		s1.Update(ctx, repository.Metrics{
			ID:    strconv.Itoa(i),
			MType: strconv.Itoa(i),
			Value: &val1,
			Delta: &delta,
		})

		s2.Update(ctx, repository.Metrics{
			ID:    strconv.Itoa(i),
			MType: strconv.Itoa(i),
			Value: &val1,
			Delta: &delta,
		})
	}

	b.ResetTimer()

	b.Run("get for slice storage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := s1.Get(ctx, "300")
			if err != nil {
				b.Failed()
			}
		}
	})
	b.Run("get for map storage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := s2.Get(ctx, "300")
			if err != nil {
				b.Failed()
			}
		}
	})
}

func BenchmarkBulkUpdate(b *testing.B) {

	s1 := inmemory.New()
	s2 := inmemory.NewFast()

	ctx := context.Background()

	// write mock items to storage
	var delta int64 = 1
	val1 := 100.1200
	bulkMetrics := make([]repository.Metrics, 0)
	for i := 0; i < 30; i++ {
		bulkMetrics = append(bulkMetrics, repository.Metrics{
			ID:    strconv.Itoa(i),
			MType: strconv.Itoa(i),
			Value: &val1,
			Delta: &delta,
		})
	}

	b.ResetTimer()

	b.Run("BulkUpdate for slice storage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := s1.BulkUpdate(ctx, bulkMetrics)
			if err != nil {
				b.Failed()
			}
		}
	})
	b.Run("BulkUpdate for map storage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := s2.BulkUpdate(ctx, bulkMetrics)
			if err != nil {
				b.Failed()
			}
		}
	})
}
