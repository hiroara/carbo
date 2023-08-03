package pipe_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/taskfn"
)

func TestTake(t *testing.T) {
	t.Parallel()

	take := taskfn.SliceToSlice(pipe.Take[string](2).AsTask())

	t.Run("MoreItems", func(t *testing.T) {
		t.Parallel()

		result, err := take(context.Background(), []string{"item1", "item2", "item3"})
		require.NoError(t, err)
		assert.Equal(t, []string{"item1", "item2"}, result)
	})

	t.Run("Equal", func(t *testing.T) {
		t.Parallel()

		result, err := take(context.Background(), []string{"item1", "item2"})
		require.NoError(t, err)
		assert.Equal(t, []string{"item1", "item2"}, result)
	})

	t.Run("LessItems", func(t *testing.T) {
		t.Parallel()

		result, err := take(context.Background(), []string{"item1"})
		require.NoError(t, err)
		assert.Equal(t, []string{"item1"}, result)
	})
}
