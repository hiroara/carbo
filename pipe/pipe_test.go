package pipe_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/pipe"
)

func double(s string) string {
	return s + s
}

func createPipeFn(fn func(string) string) (pipe.PipeFn[string, string], *bool) {
	called := false
	pipeFn := func(ctx context.Context, in <-chan string, out chan<- string) error {
		called = true
		for i := range in {
			out <- fn(i)
		}
		return nil
	}
	return pipeFn, &called
}

func TestPipeRun(t *testing.T) {
	t.Parallel()

	pipeFn, _ := createPipeFn(double)
	pipe := pipe.FromFn(pipeFn)

	in := make(chan string, 2)
	out := make(chan string, 2)
	in <- "item1"
	in <- "item2"
	close(in)

	err := pipe.Run(context.Background(), in, out)
	require.NoError(t, err)

	assert.Equal(t, "item1item1", <-out)
	assert.Equal(t, "item2item2", <-out)
}

func TestConcurrentPipe(t *testing.T) {
	t.Parallel()

	assertConcurrentPipe := func(pipe *pipe.Pipe[string, string]) {
		in := make(chan string, 2)
		out := make(chan string, 2)
		in <- "item1"
		in <- "item2"
		close(in)

		err := pipe.Run(context.Background(), in, out)
		require.NoError(t, err)

		outputs := make([]string, 0)
		for item := range out {
			outputs = append(outputs, item)
		}
		assert.ElementsMatch(t, []string{"item1item1", "item2item2"}, outputs)
	}

	t.Run("Concurrent", func(t *testing.T) {
		t.Parallel()

		pipeFn1, called1 := createPipeFn(double)
		pipeFn2, called2 := createPipeFn(double)
		pipe := pipe.Concurrent([]*pipe.Pipe[string, string]{
			pipe.FromFn(pipeFn1),
			pipe.FromFn(pipeFn2),
		})

		assertConcurrentPipe(pipe)

		assert.True(t, *called1)
		assert.True(t, *called2)
	})

	t.Run("ConcurrentFromFn", func(t *testing.T) {
		t.Parallel()

		fn, _ := createPipeFn(double)
		pipe := pipe.ConcurrentFromFn(fn, 2)

		assertConcurrentPipe(pipe)
	})
}
