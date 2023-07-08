package task_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/task"
)

func TestPipeRun(t *testing.T) {
	pipeFn := func(ctx context.Context, in <-chan string, out chan<- string) error {
		for i := range in {
			out <- i + i
		}
		return nil
	}
	pipe := task.PipeFromFn(pipeFn)

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
