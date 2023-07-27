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

func TestToSlice(t *testing.T) {
	t.Parallel()

	src := source.FromSlice([]string{"item1", "item2"})

	s := make([]string, 0)
	toSlice := sink.ToSlice(&s)

	fl := flow.FromTask(task.Connect(src.AsTask(), toSlice.AsTask(), 2))
	err := fl.Run(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []string{"item1", "item2"}, s)
}
