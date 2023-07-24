package flow

import (
	"context"

	"github.com/hiroara/carbo/config"
)

// An object to build a Flow.
type Factory interface {
	Build() (*Flow, error)
}

// A function that defines how to build a Flow.
type FactoryFn func() (*Flow, error)

// A function that defines how to build a Flow with a configuration struct.
type FactoryFnWithConfig[C any] func(cfg *C) (*Flow, error)

type factory struct {
	build FactoryFn
}

type factoryWithConfig[C any] struct {
	build   FactoryFnWithConfig[C]
	cfgPath string
}

// Create a Factory with a FactoryFn.
func NewFactory(fn FactoryFn) Factory {
	return &factory{build: fn}
}

// Create a Factory with a FactoryFnWithConfig.
// The passed cfgPath is read when building a Flow, and passed to the FactoryFnWithConfig.
func NewFactoryWithConfig[C any](fn FactoryFnWithConfig[C], cfgPath string) Factory {
	return &factoryWithConfig[C]{build: fn, cfgPath: cfgPath}
}

func (f *factory) Build() (*Flow, error) {
	return f.build()
}

func (f *factoryWithConfig[C]) Build() (*Flow, error) {
	var cfg C
	err := config.Parse(f.cfgPath, &cfg)
	if err != nil {
		return nil, err
	}

	return f.build(&cfg)
}

// Build a Flow with the passed FactoryFn and run it.
//
// This is a shorthand of creating a Factory with NewFactory, building a Flow with the Factory,
// and running the built Flow.
func Run(ctx context.Context, fn FactoryFn) error {
	f, err := NewFactory(fn).Build()
	if err != nil {
		return err
	}
	return f.Run(ctx)
}

// Build a Flow with the passed FactoryFnWithConfig and run it.
//
// This is a shorthand of creating a Factory with NewFactoryWithConfig, building a Flow with the Factory,
// and running the built Flow.
func RunWithConfig[C any](ctx context.Context, fn FactoryFnWithConfig[C], cfgPath string) error {
	f, err := NewFactoryWithConfig(fn, cfgPath).Build()
	if err != nil {
		return err
	}
	return f.Run(ctx)
}
