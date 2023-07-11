package runner_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/internal/testutils"
	"github.com/hiroara/carbo/runner"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
)

func TestRunnerRun(t *testing.T) {
	r := runner.New[testutils.Config]()

	src := source.FromSlice([]string{"item1", "item2"})

	out := make([]string, 0)
	sink := sink.ToSlice(&out)
	conn := task.Connect(src.AsTask(), sink.AsTask(), 2)

	var flowCfg *testutils.Config

	r.Define("flow1", func(cfg *testutils.Config) *flow.Flow {
		flowCfg = cfg
		return flow.FromTask(conn)
	})

	assert.Nil(t, flowCfg)

	err := r.Run(context.Background(), "flow1", "../testdata/config.yaml")
	require.NoError(t, err)

	if assert.NotNil(t, flowCfg) {
		assert.Equal(t, "thisisstring", flowCfg.StringField)
	}

	assert.Equal(t, []string{"item1", "item2"}, out)
}
