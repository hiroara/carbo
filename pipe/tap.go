package pipe

import (
	"context"

	"github.com/hiroara/carbo/task"
)

type TapOp[S any] struct {
	run TapFn[S]
}

type TapFn[S any] func(context.Context, S) error

func Tap[S any](fn TapFn[S]) *TapOp[S] {
	return &TapOp[S]{run: fn}
}

func (op *TapOp[S]) AsTask() task.Task[S, S] {
	return Map(func(ctx context.Context, el S) (S, error) {
		var zero S
		err := op.run(ctx, el)
		if err != nil {
			return zero, err
		}
		return el, nil
	}).AsTask()
}
