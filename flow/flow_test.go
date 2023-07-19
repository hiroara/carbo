package flow_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
)

func TestFlowRun(t *testing.T) {
	t.Parallel()

	src := source.FromSlice([]string{"item1", "item2"})

	items := make([]string, 0)
	sink := sink.ElementWise(func(ctx context.Context, str string) error {
		items = append(items, str)
		return nil
	})

	flow := flow.FromTask(task.Connect(src.AsTask(), sink.AsTask(), 1))
	err := flow.Run(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []string{"item1", "item2"}, items)
}
