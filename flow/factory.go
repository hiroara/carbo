package flow

import (
	"context"

	"github.com/hiroara/carbo/config"
)

type Factory interface {
	Start(ctx context.Context) error
}

type FactoryFn func() (*Flow, error)

type FactoryFnWithConfig[C any] func(cfg *C) (*Flow, error)

type factory struct {
	build FactoryFn
}

type factoryWithConfig[C any] struct {
	build   FactoryFnWithConfig[C]
	cfgPath string
}

func NewFactory(fn FactoryFn) Factory {
	return &factory{build: fn}
}

func NewFactoryWithConfig[C any](fn FactoryFnWithConfig[C], cfgPath string) Factory {
	return &factoryWithConfig[C]{build: fn, cfgPath: cfgPath}
}

func (f *factory) Start(ctx context.Context) error {
	fl, err := f.build()
	if err != nil {
		return err
	}

	return fl.Run(ctx)
}

func (f *factoryWithConfig[C]) Start(ctx context.Context) error {
	var cfg C
	err := config.Parse(f.cfgPath, &cfg)
	if err != nil {
		return err
	}

	fl, err := f.build(&cfg)
	if err != nil {
		return err
	}

	return fl.Run(ctx)
}

func Run(ctx context.Context, fn FactoryFn) error {
	return NewFactory(fn).Start(ctx)
}

func RunWithConfig[C any](ctx context.Context, fn FactoryFnWithConfig[C], cfgPath string) error {
	return NewFactoryWithConfig(fn, cfgPath).Start(ctx)
}
