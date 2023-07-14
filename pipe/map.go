package pipe

import (
	"context"

	"github.com/hiroara/carbo/task"
)

type MapOp[S, T any] struct {
	run PipeFn[S, T]
}

type MapFn[S, T any] func(S) (T, error)

func Map[S, T any](fn MapFn[S, T]) *MapOp[S, T] {
	return &MapOp[S, T]{
		run: func(ctx context.Context, in <-chan S, out chan<- T) error {
			for i := range in {
				mapped, err := fn(i)
				if err != nil {
					return err
				}
				out <- mapped
			}
			return nil
		},
	}
}

func (op *MapOp[S, T]) AsTask() task.Task[S, T] {
	return task.Task[S, T](FromFn(op.run))
}

func (op *MapOp[S, T]) Concurrent(concurrency int) *Pipe[S, T] {
	return ConcurrentFromFn(op.run, concurrency)
}
