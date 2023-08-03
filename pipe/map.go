package pipe

import (
	"context"

	"github.com/hiroara/carbo/cache"
	"github.com/hiroara/carbo/task"
)

// A Pipe task that processes an element and emits a corresponding output.
type MapOp[S, T any] struct {
	operator[S, T]
	concurrency[S, T]
	fn MapFn[S, T]
}

// A function that defines the behavior of a map operator.
type MapFn[S, T any] func(context.Context, S) (T, error)

// Create a map operator from a MapFn.
func Map[S, T any](fn MapFn[S, T]) *MapOp[S, T] {
	op := &MapOp[S, T]{fn: fn}
	op.operator.run = op.run
	op.concurrency.run = op.run
	return op
}

func (op *MapOp[S, T]) run(ctx context.Context, in <-chan S, out chan<- T) error {
	for i := range in {
		mapped, err := op.fn(ctx, i)
		if err != nil {
			return err
		}
		if err := task.Emit(ctx, out, mapped); err != nil {
			return err
		}
	}
	return nil
}

// Create a map operator with cache.
// The caching behavior is defined by the provided cache.Spec.
func MapWithCache[S, T, K, V any](fn MapFn[S, T], sp cache.Spec[S, T, K, V]) *MapOp[S, T] {
	return Map(func(ctx context.Context, el S) (T, error) {
		return cache.Run(ctx, sp, el, cache.CacheableFn[S, T](fn))
	})
}
