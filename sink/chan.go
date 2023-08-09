package sink

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A Sink task that sends elements fed by its upstream task to the passed chan.
type ToChanOp[S any] struct {
	ElementWiseOp[S]
	c chan<- S
}

// Create a ToChanOp from a chan.
func ToChan[S any](c chan<- S) *ToChanOp[S] {
	op := &ToChanOp[S]{
		ElementWiseOp: *ElementWise(func(ctx context.Context, el S) error {
			return task.Emit(ctx, c, el)
		}),
		c: c,
	}
	return op
}

// Convert this operation as a Sink.
func (op *ToChanOp[S]) AsSink(opts ...task.Option) Sink[S] {
	sin := op.ElementWiseOp.AsSink(opts...)
	sin.Defer(func() { close(op.c) })
	return sin
}

// Convert this operation as a Task.
func (op *ToChanOp[S]) AsTask(opts ...task.Option) task.Task[S, struct{}] {
	return op.AsSink(opts...)
}
