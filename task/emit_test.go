package task_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/task"
)

func TestEmit(t *testing.T) {
	t.Run("NormalCase", func(t *testing.T) {
		t.Parallel()

		out := make(chan string, 1)
		ctx := context.Background()
		err := task.Emit(ctx, out, "test")
		close(out)
		require.NoError(t, err)

		sl := make([]string, 0)
		for el := range out {
			sl = append(sl, el)
		}

		if assert.Len(t, sl, 1) {
			assert.Equal(t, "test", sl[0])
		}
	})

	t.Run("ContextCancelCase", func(t *testing.T) {
		t.Parallel()

		out := make(chan string)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Context has already been canceled.

		err := task.Emit(ctx, out, "test")
		close(out)
		require.ErrorIs(t, err, context.Canceled)
	})
}
