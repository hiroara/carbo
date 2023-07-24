package pipe

import (
	"context"

	"github.com/hiroara/carbo/cache"
	"github.com/hiroara/carbo/task"
)

// A Pipe task that processes an element and emits a corresponding output.
type MapOp[S, T any] struct {
	run MapFn[S, T]
}

// A function that defines the behavior of a map operator.
type MapFn[S, T any] func(context.Context, S) (T, error)

// Create a map operator from a MapFn.
func Map[S, T any](fn MapFn[S, T]) *MapOp[S, T] {
	return &MapOp[S, T]{
		run: fn,
	}
}

// Convert the map operator as a task.
func (op *MapOp[S, T]) AsTask() task.Task[S, T] {
	return task.Task[S, T](FromFn(op.pipeFn()))
}

// Create a concurrent Pipe from multiple map operators that have the same behavior.
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

// Create a map operator with cache.
// The caching behavior is defined by the provided cache.Spec.
func MapWithCache[S, T, K, V any](fn MapFn[S, T], sp cache.Spec[S, T, K, V]) *MapOp[S, T] {
	return Map(func(ctx context.Context, el S) (T, error) {
		return cache.Run(ctx, sp, el, cache.CacheableFn[S, T](fn))
	})
}
