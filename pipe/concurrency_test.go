package pipe_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/internal/testutils"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/task"
	"github.com/hiroara/carbo/taskfn"
)

func assertConcurrentPipe(t *testing.T, p pipe.Pipe[string, string]) {
	ctx := context.Background()
	out, err := taskfn.SliceToSlice(p.AsTask())(ctx, []string{"item1", "item2"})
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"item1item1", "item2item2"}, out)
}

func TestConcurrent(t *testing.T) {
	t.Parallel()

	t.Run("NormalCase", func(t *testing.T) {
		t.Parallel()

		pipeFn1, called1 := createPipeFn(double)
		pipeFn2, called2 := createPipeFn(double)
		p := pipe.Concurrent([]pipe.Pipe[string, string]{
			pipe.FromFn(pipeFn1),
			pipe.FromFn(pipeFn2),
		})

		assertConcurrentPipe(t, p)
		close(called1)
		close(called2)

		assert.Len(t, testutils.ReadItems(called1), 1)
		assert.Len(t, testutils.ReadItems(called2), 1)
	})

	t.Run("ErrorCase", func(t *testing.T) {
		t.Parallel()

		errFromTask := errors.New("test error")

		pipeFn1, called1 := createPipeFn(func(ctx context.Context, s string) (string, error) {
			return "", errFromTask
		})
		pipeFn2, called2 := createPipeFn(func(ctx context.Context, s string) (string, error) {
			return "", errFromTask
		})

		defer close(called1)
		defer close(called2)

		p := pipe.Concurrent([]pipe.Pipe[string, string]{
			pipe.FromFn(pipeFn1, task.WithName("pipeFn1")),
			pipe.FromFn(pipeFn2, task.WithName("pipeFn2")),
		})

		_, err := taskfn.SliceToSlice(p.AsTask())(context.Background(), []string{"item1", "item2"})
		assert.ErrorIs(t, err, errFromTask)
	})

	t.Run("NoConcurrentPipesCase", func(t *testing.T) {
		t.Parallel()

		assert.PanicsWithValue(
			t,
			"at least 1 concurrent pipe is required",
			func() { pipe.Concurrent([]pipe.Pipe[string, string]{}) },
		)
	})
}

func TestConcurrentFromFn(t *testing.T) {
	t.Run("NormalCase", func(t *testing.T) {
		t.Parallel()

		fn, called := createPipeFn(double)
		p := pipe.ConcurrentFromFn(fn, 2)

		assertConcurrentPipe(t, p)
		close(called)
		assert.Len(t, testutils.ReadItems(called), 2)
	})

	t.Run("ZeroConcurrencyCase", func(t *testing.T) {
		t.Parallel()

		fn, _ := createPipeFn(double)

		assert.PanicsWithValue(
			t,
			"at least 1 concurrent pipe is required",
			func() { pipe.ConcurrentFromFn(fn, 0) },
		)
	})

	t.Run("NegativeConcurrencyCase", func(t *testing.T) {
		t.Parallel()

		fn, _ := createPipeFn(double)

		assert.PanicsWithValue(
			t,
			"at least 1 concurrent pipe is required",
			func() { pipe.ConcurrentFromFn(fn, -1) },
		)
	})
}
