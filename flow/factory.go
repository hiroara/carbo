package flow

import (
	"context"

	"github.com/hiroara/carbo/config"
)

type FactoryFn[C any] func(cfg *C) (*Flow, error)

type Factory[C any] struct {
	build FactoryFn[C]
}

func NewFactory[C any](fn FactoryFn[C]) *Factory[C] {
	return &Factory[C]{build: fn}
}

func (f *Factory[C]) Start(ctx context.Context, cfgPath string) error {
	var cfg C
	err := config.Parse(cfgPath, &cfg)
	if err != nil {
		return err
	}

	fl, err := f.build(&cfg)
	if err != nil {
		return err
	}

	return fl.Run(ctx)
}

func Run[C any](ctx context.Context, fn FactoryFn[C], configPath string) error {
	return NewFactory(fn).Start(ctx, configPath)
}
