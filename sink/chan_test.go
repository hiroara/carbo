package sink_test

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

func TestToChan(t *testing.T) {
	t.Parallel()

	items := []string{"item1", "item2"}

	src := source.FromSlice(items)

	c := make(chan string, len(items))
	toSlice := sink.ToChan(c)

	fl := flow.FromTask(task.Connect(src.AsTask(), toSlice.AsTask(), 2))
	err := fl.Run(context.Background())
	require.NoError(t, err)

	out := make([]string, 0, len(items))
	for i := range c {
		out = append(out, i)
	}

	assert.Equal(t, items, out)
}
