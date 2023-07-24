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

func createPipeFn(fn func(string) string) (pipe.PipeFn[string, string], chan struct{}) {
	called := make(chan struct{}, 2)
	pipeFn := func(ctx context.Context, in <-chan string, out chan<- string) error {
		called <- struct{}{}
		for i := range in {
			out <- fn(i)
		}
		return nil
	}
	return pipeFn, called
}

func TestPipeRun(t *testing.T) {
	t.Parallel()

	pipeFn, called := createPipeFn(double)
	p := pipe.FromFn(pipeFn)

	deferredCalled := false
	p.Defer(func() { deferredCalled = true })

	in := make(chan string, 2)
	out := make(chan string, 2)
	in <- "item1"
	in <- "item2"
	close(in)

	err := p.Run(context.Background(), in, out)
	require.NoError(t, err)
	close(called)

	assert.Equal(t, "item1item1", <-out)
	assert.Equal(t, "item2item2", <-out)

	assert.True(t, deferredCalled)
}
