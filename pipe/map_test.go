package pipe_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/cache"
	"github.com/hiroara/carbo/cache/store"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/task"
	"github.com/hiroara/carbo/taskfn"
)

func TestMap(t *testing.T) {
	t.Parallel()

	els := []string{"item1", "item2", "item2"}

	fn := func(ctx context.Context, s string) (string, error) {
		return s + s, nil
	}
	m := pipe.Map(fn)

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

		m := pipe.Map(fn)
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

		out, err := runFlowWithMap(pipe.MapWithCache(fn, sp).AsTask())
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
