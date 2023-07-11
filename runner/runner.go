package runner

import (
	"context"

	"github.com/hiroara/carbo/flow"
)

type Runner[C any] struct {
	factories map[string]*flow.Factory[C]
}

func New[C any]() *Runner[C] {
	return &Runner[C]{factories: make(map[string]*flow.Factory[C])}
}

func (r *Runner[C]) Run(ctx context.Context, command, configPath string) error {
	return r.factories[command].Start(ctx, configPath)
}

func (r *Runner[C]) Define(name string, fn flow.FactoryFn[C]) {
	r.factories[name] = flow.NewFactory(fn)
}
