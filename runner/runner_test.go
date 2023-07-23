package runner_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/runner"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
)

func TestRunnerRun(t *testing.T) {
	t.Parallel()

	r := runner.New()

	src := source.FromSlice([]string{"item1", "item2"})

	out := make([]string, 0)
	sink := sink.ToSlice(&out)
	conn := task.Connect(src.AsTask(), sink.AsTask(), 2)
	called := false

	r.Define("flow1", flow.NewFactory(func() (*flow.Flow, error) {
		called = true
		return flow.FromTask(conn), nil
	}))

	err := r.Run(context.Background(), "flow1")
	require.NoError(t, err)

	if assert.True(t, called) {
		assert.Equal(t, []string{"item1", "item2"}, out)
	}
}
