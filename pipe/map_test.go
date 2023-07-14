package pipe_test

import (
	"context"
	"errors"
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
	t.Parallel()

	src := source.FromSlice([]string{"item1", "item2"})

	m := pipe.Map(func(s string) (string, error) {
		return s + s, nil
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

	t.Run("ErrorCase", func(t *testing.T) {
		t.Parallel()

		m := pipe.Map(func(s string) (string, error) {
			return "", errors.New("error case")
		})
		_, err := runFlowWithMap(m.AsTask())
		assert.Error(t, err)
	})

	t.Run("NoConcurrency", func(t *testing.T) {
		t.Parallel()

		out, err := runFlowWithMap(m.AsTask())
		require.NoError(t, err)

		assert.Equal(t, []string{"item1item1", "item2item2"}, out)
	})

	t.Run("Concurrent", func(t *testing.T) {
		t.Parallel()

		out, err := runFlowWithMap(m.Concurrent(2).AsTask())
		require.NoError(t, err)

		assert.ElementsMatch(t, []string{"item1item1", "item2item2"}, out)
	})
}
