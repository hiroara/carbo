package pipe

import (
	"context"
)

// A Pipe task that receives an element and emits the same element without any processing.
// This can be used to make a side effect with an input element, for example, logging elements for debug.
type TapOp[S any] struct {
	MapOp[S, S]
	operator[S, S]
	concurrency[S, S]
}

// A function that defines the behavior of a tap operator.
type TapFn[S any] func(context.Context, S) error

// Create a tap operator from a TapFn.
func Tap[S any](fn TapFn[S]) *TapOp[S] {
	op := &TapOp[S]{
		MapOp: *Map(func(ctx context.Context, el S) (S, error) {
			var zero S
			err := fn(ctx, el)
			if err != nil {
				return zero, err
			}
			return el, nil
		}),
	}
	op.operator.run = op.MapOp.run
	op.concurrency.run = op.MapOp.run
	return op
}
