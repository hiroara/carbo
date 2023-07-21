package sink

import (
	"context"

	"github.com/hiroara/carbo/task"
)

type ElementWiseOp[S any] struct {
	run SinkFn[S]
}

type ElementWiseFn[S any] func(context.Context, S) error

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

func (op *ElementWiseOp[S]) AsSink() Sink[S] {
	return FromFn(op.run)
}

func (op *ElementWiseOp[S]) AsTask() task.Task[S, struct{}] {
	return op.AsSink()
}

func (op *ElementWiseOp[S]) Concurrent(concurrency int) Sink[S] {
	return ConcurrentFromFn(op.run, concurrency)
}
