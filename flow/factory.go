package flow

import (
	"context"

	"github.com/hiroara/carbo/config"
)

type FactoryFn[C any] func(cfg *C) *Flow

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

	return f.build(&cfg).Run(ctx)
}
