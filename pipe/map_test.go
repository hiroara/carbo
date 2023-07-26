package pipe_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/cache"
	"github.com/hiroara/carbo/cache/store"
	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
)

func TestMap(t *testing.T) {
	t.Parallel()

	src := source.FromSlice([]string{"item1", "item2", "item2"})

	fn := func(ctx context.Context, s string) (string, error) {
		return s + s, nil
	}
	m := pipe.Map(fn)

	runFlowWithMap := func(mappingTask task.Task[string, string]) ([]string, error) {
		out := make([]string, 0)
		sin := sink.ToSlice(&out)

		mapped := task.Connect(src.AsTask(), mappingTask, 0)
		toSlice := task.Connect(mapped, sin.AsTask(), 2)

		err := flow.FromTask(toSlice).Run(context.Background())
		if err != nil {
			return nil, err
		}

		return out, nil
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
