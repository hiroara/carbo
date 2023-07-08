package task_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/task"
)

func TestSinkRun(t *testing.T) {
	items := make([]string, 0)
	sinkFn := func(ctx context.Context, in <-chan string) error {
		for i := range in {
			items = append(items, i)
		}
		return nil
	}
	sink := task.SinkFromFn(sinkFn)

	in := make(chan string, 2)
	out := make(chan struct{}, 1)
	in <- "item1"
	in <- "item2"
	close(in)

	err := sink.Run(context.Background(), in, out)
	require.NoError(t, err)

	assert.Equal(t, []string{"item1", "item2"}, items)
}
