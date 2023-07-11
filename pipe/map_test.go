package pipe_test

import (
	"context"
	"testing"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	src := source.FromSlice([]string{"item1", "item2"})

	m := pipe.Map(func(s string) string {
		return s + s
	})

	runFlowWithMap := func(mappingTask task.Task[string, string]) ([]string, error) {
		out := make([]string, 0)
		sin := sink.ToSlice(&out)

		mapped := task.Connect(src.AsTask(), mappingTask, 2)
		toSlice := task.Connect(mapped, sin.AsTask(), 2)

		err := flow.FromTask(toSlice).Run(context.Background())
		if err != nil {
			return nil, err
		}

		return out, nil
	}

	t.Run("NoConcurrency", func(t *testing.T) {
		out, err := runFlowWithMap(m.AsTask())
		require.NoError(t, err)

		assert.Equal(t, []string{"item1item1", "item2item2"}, out)
	})

	t.Run("Concurrent", func(t *testing.T) {
		out, err := runFlowWithMap(m.Concurrent(2).AsTask())
		require.NoError(t, err)

		assert.ElementsMatch(t, []string{"item1item1", "item2item2"}, out)
	})
}
