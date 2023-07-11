package flow_test

import (
	"context"
	"testing"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/internal/testutils"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Config struct {
	Value string `yaml:"value"`
}

func TestFactoryStart(t *testing.T) {
	t.Parallel()

	var flowCfg *testutils.Config

	fac := flow.NewFactory(func(cfg *testutils.Config) *flow.Flow {
		flowCfg = cfg
		src := source.FromSlice([]string{"item1", "item2"})

		items := make([]string, 0)
		sink := sink.ElementWise(func(str string) error {
			items = append(items, str)
			return nil
		})
		return flow.FromTask(task.Connect(src.AsTask(), sink.AsTask(), 1))
	})

	assert.Nil(t, flowCfg)

	err := fac.Start(context.Background(), "../testdata/config.yaml")
	require.NoError(t, err)

	assert.Equal(t, "thisisstring", flowCfg.StringField) // Decoded config is passed to the factory function
}
