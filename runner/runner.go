package runner

import (
	"context"

	"github.com/hiroara/carbo/flow"
)

type Runner struct {
	factories map[string]flow.Factory
}

func New() *Runner {
	return &Runner{factories: make(map[string]flow.Factory)}
}

func (r *Runner) Run(ctx context.Context, command string) error {
	return r.factories[command].Start(ctx)
}

func (r *Runner) Define(name string, factory flow.Factory) {
	r.factories[name] = factory
}
