package task_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/task"
)

func TestConnectionRun(t *testing.T) {
	src := task.SourceFromFn(func(ctx context.Context, out chan<- string) error {
		out <- "item1"
		out <- "item2"
		return nil
	})

	items := make([]string, 0)
	sink := task.SinkFromFn(func(ctx context.Context, in <-chan string) error {
		for i := range in {
			items = append(items, i)
		}
		return nil
	})

	conn := task.Connect(src.AsTask(), sink.AsTask(), 1)

	in := make(chan struct{})
	out := make(chan struct{})
	close(in)

	err := conn.Run(context.Background(), in, out)
	require.NoError(t, err)

	assert.Equal(t, []string{"item1", "item2"}, items)
}
