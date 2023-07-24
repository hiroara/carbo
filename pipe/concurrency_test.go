package pipe_test

import (
	"context"
	"testing"

	"github.com/hiroara/carbo/internal/testutils"
	"github.com/hiroara/carbo/pipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConcurrentPipe(t *testing.T) {
	t.Parallel()

	assertConcurrentPipe := func(p *pipe.Pipe[string, string]) {
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

	t.Run("Concurrent", func(t *testing.T) {
		t.Parallel()

		pipeFn1, called1 := createPipeFn(double)
		pipeFn2, called2 := createPipeFn(double)
		p := pipe.Concurrent([]*pipe.Pipe[string, string]{
			pipe.FromFn(pipeFn1),
			pipe.FromFn(pipeFn2),
		})

		assertConcurrentPipe(p)
		close(called1)
		close(called2)

		assert.Len(t, testutils.ReadItems(called1), 1)
		assert.Len(t, testutils.ReadItems(called2), 1)
	})

	t.Run("ConcurrentFromFn", func(t *testing.T) {
		t.Parallel()

		fn, called := createPipeFn(double)
		p := pipe.ConcurrentFromFn(fn, 2)

		assertConcurrentPipe(p)
		close(called)
		assert.Len(t, testutils.ReadItems(called), 2)
	})
}
