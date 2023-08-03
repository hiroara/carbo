package pipe_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/internal/testutils"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/taskfn"
)

func TestTap(t *testing.T) {
	t.Parallel()

	t.Run("NormalCase", func(t *testing.T) {
		t.Parallel()

		called := make(chan string, 2)
		tap := taskfn.SliceToSlice(pipe.Tap(func(ctx context.Context, el string) error {
			called <- el
			return nil
		}).AsTask())

		out, err := tap(context.Background(), []string{"item1", "item2"})
		require.NoError(t, err)
		close(called)

		assert.Equal(t, []string{"item1", "item2"}, testutils.ReadItems(called))
		assert.Equal(t, []string{"item1", "item2"}, out)
	})

	t.Run("ErrorCase", func(t *testing.T) {
		t.Parallel()

		tapErr := errors.New("test error")
		tap := taskfn.SliceToSlice(pipe.Tap(func(ctx context.Context, el string) error {
			return tapErr
		}).AsTask())

		_, err := tap(context.Background(), []string{"item1", "item2"})
		require.ErrorIs(t, err, tapErr)
	})

	t.Run("Concurrent", func(t *testing.T) {
		t.Parallel()

		called := make(chan string, 2)
		tap := taskfn.SliceToSlice(pipe.Tap(func(ctx context.Context, el string) error {
			called <- el
			return nil
		}).Concurrent(2).AsTask())

		out, err := tap(context.Background(), []string{"item1", "item2"})
		require.NoError(t, err)
		close(called)

		assert.ElementsMatch(t, []string{"item1", "item2"}, testutils.ReadItems(called))
		assert.ElementsMatch(t, []string{"item1", "item2"}, out)
	})
}
