package pipe_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/task"
)

func double(ctx context.Context, s string) (string, error) {
	return s + s, nil
}

func createPipeFn(fn func(context.Context, string) (string, error)) (pipe.PipeFn[string, string], chan struct{}) {
	called := make(chan struct{}, 2)
	pipeFn := func(ctx context.Context, in <-chan string, out chan<- string) error {
		called <- struct{}{}
		for i := range in {
			el, err := fn(ctx, i)
			if err != nil {
				return err
			}
			if err := task.Emit(ctx, out, el); err != nil {
				return err
			}
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
