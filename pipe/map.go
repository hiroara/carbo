package pipe

import (
	"context"

	"github.com/hiroara/carbo/cache"
	"github.com/hiroara/carbo/task"
)

type MapOp[S, T any] struct {
	run MapFn[S, T]
}

type MapFn[S, T any] func(context.Context, S) (T, error)

func Map[S, T any](fn MapFn[S, T]) *MapOp[S, T] {
	return &MapOp[S, T]{
		run: fn,
	}
}

func (op *MapOp[S, T]) AsTask() task.Task[S, T] {
	return task.Task[S, T](FromFn(op.pipeFn()))
}

func (op *MapOp[S, T]) Concurrent(concurrency int) *Pipe[S, T] {
	return ConcurrentFromFn(op.pipeFn(), concurrency)
}

func (op *MapOp[S, T]) pipeFn() PipeFn[S, T] {
	return func(ctx context.Context, in <-chan S, out chan<- T) error {
		for i := range in {
			mapped, err := op.run(ctx, i)
			if err != nil {
				return err
			}
			out <- mapped
		}
		return nil
	}
}

func MapWithCache[S, T, K, V any](fn MapFn[S, T], sp cache.Spec[S, T, K, V]) *MapOp[S, T] {
	return Map(func(ctx context.Context, el S) (T, error) {
		var zero T

		ent, err := cache.GetEntry(sp, el)
		if err != nil {
			return zero, err
		}

		return ent.Run(ctx, cache.CacheableFn[S, T](fn))
	})
}
