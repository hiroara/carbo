package sink

import (
	"context"
)

// A Sink task that processes elements fed by its upstream task one by one.
type ElementWiseOp[S any] struct {
	operator[S]
	fn ElementWiseFn[S]
}

// A function that defines the behavior of an elementwise operator.
type ElementWiseFn[S any] func(context.Context, S) error

// Create an elementwise operator from an ElementWiseFn.
func ElementWise[S any](fn ElementWiseFn[S]) *ElementWiseOp[S] {
	op := &ElementWiseOp[S]{fn: fn}
	op.operator.run = op.run
	return op
}

func (op *ElementWiseOp[S]) run(ctx context.Context, in <-chan S) error {
	for el := range in {
		if err := op.fn(ctx, el); err != nil {
			return err
		}
	}
	return nil
}

// Create a concurrent Sink from multiple elementwise operators that have the same behavior.
func (op *ElementWiseOp[S]) Concurrent(concurrency int) Sink[S] {
	return ConcurrentFromFn(op.run, concurrency)
}
