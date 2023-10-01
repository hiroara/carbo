package pipe

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

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

// Create a concurrent Pipe to apply the map operator.
//
// Unlike a Pipe created with Concurrent, a concurrent Pipe created with this ConcurrentPreservingOrder,
// preserves the order of elements.
func (op *MapOp[S, T]) ConcurrentPreservingOrder(concurrency int, opts ...task.Option) Pipe[S, T] {
	if concurrency <= 0 {
		panic("at least 1 concurrency is required")
	}

	return FromFn(func(ctx context.Context, in <-chan S, out chan<- T) error {
		grp, ctx := errgroup.WithContext(ctx)
		ins, outs, agg := duplicateOutChanPreservingOrder(in, out, concurrency)
		for idx := 0; idx < concurrency; idx++ {
			i := ins[idx]
			o := outs[idx]
			grp.Go(func() error {
				defer close(o)
				return op.run(ctx, i, o)
			})
		}
		grp.Go(func() error { return agg(ctx) })
		return grp.Wait()
	}, opts...)
}

func duplicateOutChanPreservingOrder[S, T any](
	in <-chan S, out chan<- T, n int,
) ([]<-chan S, []chan<- T, func(context.Context) error) {
	if n <= 0 {
		panic(fmt.Sprintf("argument n must be a positive value but received %d", n))
	}

	ins := make([]chan S, n)
	outs := make([]chan T, n)
	insRet := make([]<-chan S, n)
	outsRet := make([]chan<- T, n)
	for idx := 0; idx < n; idx++ {
		i := make(chan S)
		o := make(chan T)
		ins[idx] = i
		insRet[idx] = i
		outs[idx] = o
		outsRet[idx] = o
	}
	return insRet, outsRet, func(ctx context.Context) error {
		grp, ctx := errgroup.WithContext(ctx)

		grp.Go(func() error {
			defer func() {
				for _, i := range ins {
					close(i)
				}
			}()
			idx := 0
			for el := range in {
				ins[idx] <- el
				idx = (idx + 1) % n
			}
			return nil
		})

		grp.Go(func() error {
			for {
				for _, o := range outs {
					el, ok := <-o
					if !ok {
						return nil
					}
					if err := task.Emit(ctx, out, el); err != nil {
						return err
					}
				}
			}
		})

		return grp.Wait()
	}
}
