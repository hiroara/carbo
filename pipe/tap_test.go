package pipe_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/internal/testutils"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
)

func TestTap(t *testing.T) {
	t.Parallel()

	src := source.FromSlice([]string{"item1", "item2"})

	t.Run("NormalCase", func(t *testing.T) {
		called := make(chan string, 2)
		tap := task.Connect(src.AsTask(), pipe.Tap(func(ctx context.Context, el string) error {
			called <- el
			return nil
		}).AsTask(), 1)

		out := []string{}
		err := flow.FromTask(task.Connect(
			tap,
			sink.ToSlice(&out).AsTask(),
			1,
		)).Run(context.Background())
		require.NoError(t, err)
		close(called)

		assert.Equal(t, []string{"item1", "item2"}, testutils.ReadItems(called))
		assert.Equal(t, []string{"item1", "item2"}, out)
	})

	t.Run("ErrorCase", func(t *testing.T) {
		tapErr := errors.New("test error")
		tap := task.Connect(src.AsTask(), pipe.Tap(func(ctx context.Context, el string) error {
			return tapErr
		}).AsTask(), 1)

		out := []string{}
		err := flow.FromTask(task.Connect(
			tap,
			sink.ToSlice(&out).AsTask(),
			1,
		)).Run(context.Background())
		require.ErrorIs(t, err, tapErr)
	})
}
