package flow_test

import (
	"context"
	"errors"
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

func createFactoryFn() (flow.FactoryFn[testutils.Config], *testutils.Config) {
	var flowCfg testutils.Config

	fn := func(cfg *testutils.Config) (*flow.Flow, error) {
		flowCfg = *cfg
		src := source.FromSlice([]string{"item1", "item2"})

		items := make([]string, 0)
		sink := sink.ToSlice(&items)
		return flow.FromTask(task.Connect(src.AsTask(), sink.AsTask(), 1)), nil
	}

	return fn, &flowCfg
}

func TestFactoryStart(t *testing.T) {
	t.Parallel()

	t.Run("NormalCase", func(t *testing.T) {
		t.Parallel()

		fn, cfg := createFactoryFn()
		fac := flow.NewFactory(fn)

		assert.Zero(t, *cfg)

		err := fac.Start(context.Background(), "../testdata/config.yaml")
		require.NoError(t, err)

		assert.Equal(t, "thisisstring", cfg.StringField) // Decoded config is passed to the factory function
	})

	t.Run("ErrorCase", func(t *testing.T) {
		t.Parallel()

		factoryErr := errors.New("test error")

		fac := flow.NewFactory(func(cfg *struct{}) (*flow.Flow, error) {
			return nil, factoryErr
		})

		err := fac.Start(context.Background(), "../testdata/config.yaml")
		require.ErrorIs(t, err, factoryErr)
	})
}

func TestRun(t *testing.T) {
	t.Parallel()

	fn, cfg := createFactoryFn()

	assert.Zero(t, *cfg)

	err := flow.Run(context.Background(), fn, "../testdata/config.yaml")
	require.NoError(t, err)

	assert.Equal(t, "thisisstring", cfg.StringField) // Decoded config is passed to the factory function
}
