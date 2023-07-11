package sink_test

import (
	"context"
	"testing"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestElementWise(t *testing.T) {
	t.Parallel()

	src := source.FromSlice([]string{"item1", "item2"})

	items := make([]string, 0)
	sink := sink.ElementWise(func(s string) error {
		items = append(items, s)
		return nil
	})

	conn := task.Connect(src.AsTask(), sink.AsTask(), 2)
	err := flow.FromTask(conn).Run(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []string{"item1", "item2"}, items)
}
