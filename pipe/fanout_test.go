package pipe_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFanout(t *testing.T) {
	t.Parallel()

	src := source.FromSlice([]int{1, 2, 3})

	setup := func(fo *pipe.FanoutOp[int, string, []string]) *bool {
		fo.Add(pipe.Map(func(ctx context.Context, i int) (string, error) {
			return strconv.FormatInt(int64(i), 10), nil
		}).AsTask(), 2, 2)
		fo.Add(pipe.Map(func(ctx context.Context, i int) (string, error) {
			return strconv.FormatInt(int64(i*10), 10), nil
		}).AsTask(), 2, 2)

		deferredCalled := false
		fo.Defer(func() { deferredCalled = true })
		return &deferredCalled
	}

	t.Run("WithFanoutAggregateFn", func(t *testing.T) {
		t.Parallel()

		fo := pipe.Fanout[int](func(ctx context.Context, ss []string, out chan<- []string) error {
			for _, s := range ss {
				out <- []string{s}
			}
			return nil
		})

		deferredCalled := setup(fo)

		conn := task.Connect(src.AsTask(), fo.AsTask(), 2)

		out := make([][]string, 0)
		err := flow.FromTask(task.Connect(conn, sink.ToSlice(&out).AsTask(), 2)).Run(context.Background())
		require.NoError(t, err)

		assert.Equal(t, [][]string{{"1"}, {"10"}, {"2"}, {"20"}, {"3"}, {"30"}}, out)

		assert.True(t, *deferredCalled)
	})

	t.Run("WithFanoutMapFn", func(t *testing.T) {
		t.Parallel()

		fo := pipe.FanoutWithMap[int](func(ctx context.Context, ss []string) ([]string, error) {
			return ss, nil
		})

		deferredCalled := setup(fo)

		conn := task.Connect(src.AsTask(), fo.AsTask(), 2)

		out := make([][]string, 0)
		err := flow.FromTask(task.Connect(conn, sink.ToSlice(&out).AsTask(), 2)).Run(context.Background())
		require.NoError(t, err)

		assert.Equal(t, [][]string{{"1", "10"}, {"2", "20"}, {"3", "30"}}, out)

		assert.True(t, *deferredCalled)
	})
}
