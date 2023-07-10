package task_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
)

func TestFlowRun(t *testing.T) {
	t.Parallel()

	src := source.FromFn(func(ctx context.Context, out chan<- string) error {
		out <- "item1"
		out <- "item2"
		return nil
	})

	items := make([]string, 0)
	sink := sink.FromFn(func(ctx context.Context, in <-chan string) error {
		for i := range in {
			items = append(items, i)
		}
		return nil
	})

	flow := task.FlowFromTask(task.Connect(src.AsTask(), sink.AsTask(), 1))
	err := flow.Run(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []string{"item1", "item2"}, items)
}
