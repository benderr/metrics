package inmemory_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/benderr/metrics/internal/server/repository"
	"github.com/benderr/metrics/internal/server/repository/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func ExampleKeyValueMetricRepository_Get() {
	opCtx, opCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer opCancel()

	s := inmemory.NewFast()

	s.Get(opCtx, "some-test-id")
}

func ExampleKeyValueMetricRepository_BulkUpdate() {
	opCtx, opCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer opCancel()

	s := inmemory.NewFast()

	var delta int64 = 1
	value := 100.1200

	s.BulkUpdate(opCtx, []repository.Metrics{
		{
			ID:    "some-test-id",
			MType: "gauge",
			Value: &value,
		},
		{
			ID:    "some-test-id-2",
			MType: "counter",
			Delta: &delta,
		},
	})
}

func TestInMemoryMetricRepository_New(t *testing.T) {
	t.Run("storage_add", func(t *testing.T) {
		s := inmemory.New()
		ctx := context.Background()
		val := 1.1
		v, err := s.Update(ctx, repository.Metrics{ID: "test", MType: "gauge", Value: &val})

		require.NoError(t, err)
		assert.Equal(t, "test", v.ID)
		assert.Equal(t, val, *v.Value)

		item, err := s.Get(ctx, "test")

		require.NoError(t, err)
		assert.Equal(t, "test", item.ID)
		assert.Equal(t, val, *item.Value)
	})

	t.Run("storage_update", func(t *testing.T) {
		s := inmemory.New()
		ctx := context.Background()
		var val int64 = 1
		_, errAdd := s.Update(ctx, repository.Metrics{ID: "test", Delta: &val, MType: "counter"})

		require.NoError(t, errAdd)

		var val2 int64 = 3
		_, errUpdate := s.Update(ctx, repository.Metrics{ID: "test", Delta: &val2, MType: "counter"})

		require.NoError(t, errUpdate)

		item, err := s.Get(ctx, "test")

		require.NoError(t, err)
		assert.Equal(t, val+val2, *item.Delta)
	})

	t.Run("storage_get_list", func(t *testing.T) {
		s := inmemory.New()
		ctx := context.Background()
		val := 1.1
		_, err := s.Update(ctx, repository.Metrics{ID: "test", Value: &val})
		_, err2 := s.Update(ctx, repository.Metrics{ID: "test2", Value: &val})

		require.NoError(t, err)
		require.NoError(t, err2)

		items, err := s.GetList(ctx)

		require.NoError(t, err)
		assert.Equal(t, 2, len(items))
	})

	t.Run("storage_bulk_update", func(t *testing.T) {
		s := inmemory.New()
		ctx := context.Background()
		var delta int64 = 1
		val := 1.1
		err := s.BulkUpdate(ctx, []repository.Metrics{
			{ID: "test", Delta: &delta, MType: "counter"},
			{ID: "test2", Value: &val, MType: "gauge"},
		})

		require.NoError(t, err)

		items, err := s.GetList(ctx)

		require.NoError(t, err)
		assert.Equal(t, 2, len(items))
	})
}

func TestInMemoryMetricRepository_NewFast(t *testing.T) {
	t.Run("storage_add", func(t *testing.T) {
		s := inmemory.NewFast()
		ctx := context.Background()
		val := 1.1
		v, err := s.Update(ctx, repository.Metrics{ID: "test", MType: "gauge", Value: &val})

		require.NoError(t, err)
		assert.Equal(t, "test", v.ID)
		assert.Equal(t, val, *v.Value)

		item, err := s.Get(ctx, "test")

		require.NoError(t, err)
		assert.Equal(t, "test", item.ID)
		assert.Equal(t, val, *item.Value)
	})

	t.Run("storage_update", func(t *testing.T) {
		s := inmemory.NewFast()
		ctx := context.Background()
		var val int64 = 1
		_, errAdd := s.Update(ctx, repository.Metrics{ID: "test", Delta: &val, MType: "counter"})

		require.NoError(t, errAdd)

		var val2 int64 = 3
		_, errUpdate := s.Update(ctx, repository.Metrics{ID: "test", Delta: &val2, MType: "counter"})

		require.NoError(t, errUpdate)

		item, err := s.Get(ctx, "test")

		require.NoError(t, err)
		assert.Equal(t, val+val2, *item.Delta)
	})

	t.Run("storage_get_list", func(t *testing.T) {
		s := inmemory.NewFast()
		ctx := context.Background()
		val := 1.1
		_, err := s.Update(ctx, repository.Metrics{ID: "test", Value: &val})
		_, err2 := s.Update(ctx, repository.Metrics{ID: "test2", Value: &val})

		require.NoError(t, err)
		require.NoError(t, err2)

		items, err := s.GetList(ctx)

		require.NoError(t, err)
		assert.Equal(t, 2, len(items))
	})

	t.Run("storage_bulk_update", func(t *testing.T) {
		s := inmemory.NewFast()
		ctx := context.Background()
		var delta int64 = 1
		val := 1.1
		err := s.BulkUpdate(ctx, []repository.Metrics{
			{ID: "test", Delta: &delta, MType: "counter"},
			{ID: "test2", Value: &val, MType: "gauge"},
		})

		require.NoError(t, err)

		items, err := s.GetList(ctx)

		require.NoError(t, err)
		assert.Equal(t, 2, len(items))
	})
}
