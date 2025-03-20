package pipe_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/cache"
	"github.com/hiroara/carbo/cache/store"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/task"
	"github.com/hiroara/carbo/taskfn"
)

func doubleString(ctx context.Context, s string) (string, error) {
	return s + s, nil
}
func TestMap(t *testing.T) {
	t.Parallel()

	els := []string{"item1", "item2", "item2"}

	m := pipe.Map(doubleString)

	runFlowWithMap := func(mappingTask task.Task[string, string]) ([]string, error) {
		tfn := taskfn.SliceToSlice(mappingTask)
		return tfn(context.Background(), els)
	}

	t.Run("ErrorCase", func(t *testing.T) {
		t.Parallel()

		m := pipe.Map(func(ctx context.Context, s string) (string, error) {
			return "", errors.New("error case")
		})
		_, err := runFlowWithMap(m.AsTask())
		assert.Error(t, err)
	})

	t.Run("NoConcurrency", func(t *testing.T) {
		t.Parallel()

		out, err := runFlowWithMap(m.AsTask())
		require.NoError(t, err)

		assert.Equal(t, []string{"item1item1", "item2item2", "item2item2"}, out)
	})

	t.Run("Concurrent", func(t *testing.T) {
		t.Parallel()

		out, err := runFlowWithMap(m.Concurrent(2).AsTask())
		require.NoError(t, err)

		assert.ElementsMatch(t, []string{"item1item1", "item2item2", "item2item2"}, out)
	})

	t.Run("ConcurrentPreservingOrder", func(t *testing.T) {
		t.Parallel()

		m := pipe.Map(doubleString)
		tfn := taskfn.SliceToSlice(m.ConcurrentPreservingOrder(2).AsTask())
		els := []string{"item1", "item2", "item3", "item4", "item5", "item6", "item7", "item8", "item9", "item10"}
		out, err := tfn(context.Background(), els)
		require.NoError(t, err)

		expected := make([]string, len(els))
		for i, el := range els {
			expected[i] = el + el
		}

		assert.Equal(t, expected, out)
	})

	t.Run("Cache", func(t *testing.T) {
		t.Parallel()

		cs := store.NewMemoryStore[string]()
		sp := cache.NewRawSpec[string, string, string](
			cs,
			func(el string) (*cache.StoreKey[string], error) {
				return cache.Key(el), nil
			},
		)

		out, err := runFlowWithMap(pipe.MapWithCache(doubleString, sp).AsTask())
		require.NoError(t, err)

		assert.ElementsMatch(t, []string{"item1item1", "item2item2", "item2item2"}, out)

		ctx := context.Background()
		vp, err := cs.Get(ctx, "item1")
		require.NoError(t, err)

		if assert.NotNil(t, vp) {
			assert.Equal(t, "item1item1", *vp)
		}
	})
}

func doubleStringWithSleep(ctx context.Context, s string) (string, error) {
	time.Sleep(time.Duration(rand.Float64()*100) * time.Microsecond)
	return doubleString(ctx, s)
}

const mapBenchTime = 100

func BenchmarkMap(b *testing.B) {
	els := make([]string, mapBenchTime)
	for i := range els {
		els[i] = fmt.Sprintf("item%d", i)
	}

	m := pipe.Map(doubleStringWithSleep)
	tfn := taskfn.SliceToSlice(m.AsTask())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tfn(context.Background(), els)
	}
}

func BenchmarkMapConcurrent(b *testing.B) {
	els := make([]string, mapBenchTime)
	for i := range els {
		els[i] = fmt.Sprintf("item%d", i)
	}

	m := pipe.Map(doubleStringWithSleep)
	tfn := taskfn.SliceToSlice(m.Concurrent(4).AsTask())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tfn(context.Background(), els)
	}
}

func BenchmarkMapConcurrentPreservingOrder(b *testing.B) {
	els := make([]string, mapBenchTime)
	for i := range els {
		els[i] = fmt.Sprintf("item%d", i)
	}

	m := pipe.Map(doubleStringWithSleep)
	tfn := taskfn.SliceToSlice(m.ConcurrentPreservingOrder(4).AsTask())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tfn(context.Background(), els)
	}
}
