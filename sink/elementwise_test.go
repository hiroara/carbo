package sink_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/internal/testutils"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
)

func TestElementWise(t *testing.T) {
	t.Parallel()

	createAppendOp := func(sl chan<- string) *sink.ElementWiseOp[string] {
		return sink.ElementWise(func(ctx context.Context, s string) error {
			sl <- s
			return nil
		})
	}

	runFlowWithSink := func(sinkTask task.Task[string, struct{}]) error {
		src := source.FromSlice([]string{"item1", "item2"})

		conn := task.Connect(src.AsTask(), sinkTask, 2)
		err := flow.FromTask(conn).Run(context.Background())
		return err
	}

	t.Run("ErrorCase", func(t *testing.T) {
		t.Parallel()

		op := sink.ElementWise(func(ctx context.Context, s string) error {
			return errors.New("test error")
		})

		err := runFlowWithSink(op.AsTask())
		require.Error(t, err)
	})

	t.Run("NoConcurrency", func(t *testing.T) {
		t.Parallel()

		out := make(chan string, 2)
		op := createAppendOp(out)

		err := runFlowWithSink(op.AsTask())
		require.NoError(t, err)
		close(out)

		assert.Equal(t, []string{"item1", "item2"}, testutils.ReadItems(out))
	})

	t.Run("Concurrent", func(t *testing.T) {
		t.Parallel()

		out := make(chan string, 2)
		op := createAppendOp(out)

		err := runFlowWithSink(op.Concurrent(2).AsTask())
		require.NoError(t, err)
		close(out)

		assert.ElementsMatch(t, []string{"item1", "item2"}, testutils.ReadItems(out))
	})
}
