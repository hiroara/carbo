package pipe

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A Pipe task that receives an element and emits the same element without any processing.
// This can be used to make a side effect with an input element, for example, logging elements for debug.
type TapOp[S any] struct {
	run TapFn[S]
}

// A function that defines the behavior of a tap operator.
type TapFn[S any] func(context.Context, S) error

// Create a tap operator from a TapFn.
func Tap[S any](fn TapFn[S]) *TapOp[S] {
	return &TapOp[S]{run: fn}
}

// Convert the tap operator as a Pipe.
func (op *TapOp[S]) AsPipe() Pipe[S, S] {
	return Map(func(ctx context.Context, el S) (S, error) {
		var zero S
		err := op.run(ctx, el)
		if err != nil {
			return zero, err
		}
		return el, nil
	}).AsPipe()
}

// Convert the tap operator as a task.
func (op *TapOp[S]) AsTask() task.Task[S, S] {
	return op.AsPipe().AsTask()
}
