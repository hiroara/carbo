package pipe_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/internal/testutils"
	"github.com/hiroara/carbo/pipe"
)

func assertConcurrentPipe(t *testing.T, p pipe.Pipe[string, string]) {
	in := make(chan string, 2)
	out := make(chan string, 2)
	in <- "item1"
	in <- "item2"
	close(in)

	err := p.Run(context.Background(), in, out)
	require.NoError(t, err)

	outputs := make([]string, 0)
	for item := range out {
		outputs = append(outputs, item)
	}
	assert.ElementsMatch(t, []string{"item1item1", "item2item2"}, outputs)
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
