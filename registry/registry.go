package registry

import (
	"context"
	"errors"
	"fmt"

	"github.com/hiroara/carbo/flow"
)

// A object to register how flows are created.
//
// This can be used to define multiple flows within an executable and invoke one of them.
// A typical usage of a Registry is that running a registered flow by selecting it with an argument for an executable.
// This works like an executable that takes an argument as a subcommand.
type Registry struct {
	factories map[string]flow.Factory
}

// Create a Registry.
func New() *Registry {
	return &Registry{factories: make(map[string]flow.Factory)}
}

var ErrNoMatchingFlow = errors.New("no matching flow is found")

// Build a registered Flow selected with the passed name, and run it.
func (r *Registry) Run(ctx context.Context, name string) error {
	fac, ok := r.factories[name]
	if !ok {
		return fmt.Errorf("%w with name \"%s\"", ErrNoMatchingFlow, name)
	}

	f, err := fac.Build()
	if err != nil {
		return err
	}
	return f.Run(ctx)
}

// Register a Factory to build a Flow.
func (r *Registry) Register(name string, factory flow.Factory) {
	r.factories[name] = factory
}
