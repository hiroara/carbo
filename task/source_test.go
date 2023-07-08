package task_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/task"
)

func TestSourceRun(t *testing.T) {
	srcFn := func(ctx context.Context, out chan<- string) error {
		out <- "item1"
		out <- "item2"
		return nil
	}
	src := task.SourceFromFn(srcFn)

	in := make(chan struct{})
	out := make(chan string, 2)
	close(in)

	err := src.Run(context.Background(), in, out)
	require.NoError(t, err)

	assert.Equal(t, "item1", <-out)
	assert.Equal(t, "item2", <-out)
}
