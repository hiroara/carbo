package pipe_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/taskfn"
)

func TestBatch(t *testing.T) {
	t.Parallel()

	t.Run("Mod=0", func(t *testing.T) {
		t.Parallel()

		batch := taskfn.SliceToSlice(pipe.Batch[string](2).AsTask())
		result, err := batch(context.Background(), []string{"a", "b", "c", "d"})
		require.NoError(t, err)
		assert.Equal(t, [][]string{{"a", "b"}, {"c", "d"}}, result)
	})

	t.Run("Mod=1", func(t *testing.T) {
		t.Parallel()

		batch := taskfn.SliceToSlice(pipe.Batch[string](2).AsTask())
		result, err := batch(context.Background(), []string{"a", "b", "c", "d", "e"})
		require.NoError(t, err)
		assert.Equal(t, [][]string{{"a", "b"}, {"c", "d"}, {"e"}}, result)
	})
}
