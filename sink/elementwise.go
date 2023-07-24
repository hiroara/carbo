package sink

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A Sink task that processes elements fed by its upstream task one by one.
type ElementWiseOp[S any] struct {
	run SinkFn[S]
}

// A function that defines the behavior of an elementwise operator.
type ElementWiseFn[S any] func(context.Context, S) error

// Create an elementwise operator from an ElementWiseFn.
func ElementWise[S any](fn ElementWiseFn[S]) *ElementWiseOp[S] {
	return &ElementWiseOp[S]{
		run: func(ctx context.Context, in <-chan S) error {
			for el := range in {
				if err := fn(ctx, el); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

// Convert the elementwise operator as a sink.
func (op *ElementWiseOp[S]) AsSink() Sink[S] {
	return FromFn(op.run)
}

// Convert the elementwise operator as a task.
func (op *ElementWiseOp[S]) AsTask() task.Task[S, struct{}] {
	return op.AsSink()
}

// Create a concurrent Sink from multiple elementwise operators that have the same behavior.
func (op *ElementWiseOp[S]) Concurrent(concurrency int) Sink[S] {
	return ConcurrentFromFn(op.run, concurrency)
}
