package pipe

import (
	"context"
	"fmt"
	"reflect"

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

type spreadFn[S any] func(ctx context.Context, idx chan<- int, in <-chan S, ins []chan<- S) error

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
// Unlike a Pipe created with Concurrent, a concurrent Pipe created with this ConcurrentPreservingOrder
// preserves the order of elements.
func (op *MapOp[S, T]) ConcurrentPreservingOrder(concurrency int, opts ...task.Option) Pipe[S, T] {
	return op.concurrentPreservingOrder(concurrency, spread[S], opts...)
}

// Create a concurrent Pipe to apply the map operator.
//
// A concurrent Pipe created with this method preserves the order of elements like ConcurrentPreservingOrder.
// And, it additionally allows you to specify which pipe to be used for each element by providing the bucket function.
// Use pipe.ConcurrencyIndex to get the index of the current pipe.
func (op *MapOp[S, T]) StickyConcurrentPreservingOrder(
	bucket func(S) int,
	concurrency int,
	opts ...task.Option,
) Pipe[S, T] {
	return op.concurrentPreservingOrder(
		concurrency,
		func(ctx context.Context, idx chan<- int, in <-chan S, outs []chan<- S) error {
			return stickySpread(ctx, idx, in, outs, bucket)
		},
		opts...,
	)
}

func (op *MapOp[S, T]) concurrentPreservingOrder(concurrency int, spread spreadFn[S], opts ...task.Option) Pipe[S, T] {
	if concurrency <= 0 {
		panic("at least 1 concurrency is required")
	}

	return FromFn(func(ctx context.Context, in <-chan S, out chan<- T) error {
		grp, ctx := errgroup.WithContext(ctx)
		ins, outs, agg := duplicateOutChanPreservingOrder(in, out, concurrency, spread)
		for idx := 0; idx < concurrency; idx++ {
			ci := idx
			i := ins[idx]
			o := outs[idx]
			grp.Go(func() error {
				defer close(o)
				ctx := context.WithValue(ctx, concurrencyIndex, ci)
				return op.run(ctx, i, o)
			})
		}
		grp.Go(func() error { return agg(ctx) })
		return grp.Wait()
	}, opts...)
}

func duplicateOutChanPreservingOrder[S, T any](
	in <-chan S, out chan<- T, n int,
	spread spreadFn[S],
) ([]<-chan S, []chan<- T, func(context.Context) error) {
	if n <= 0 {
		panic(fmt.Sprintf("argument n must be a positive value but received %d", n))
	}

	ins := make([]chan<- S, n)
	outs := make([]<-chan T, n)
	insRet := make([]<-chan S, n)
	outsRet := make([]chan<- T, n)
	for idx := 0; idx < n; idx++ {
		i := make(chan S)
		o := make(chan T, n)
		ins[idx] = i
		insRet[idx] = i
		outs[idx] = o
		outsRet[idx] = o
	}
	return insRet, outsRet, func(ctx context.Context) error {
		grp, ctx := errgroup.WithContext(ctx)

		idx := make(chan int, n)

		grp.Go(func() error { return spread(ctx, idx, in, ins) })

		grp.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return context.Cause(ctx)
				case i := <-idx:
					el, ok := <-outs[i]
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

func spread[S any](ctx context.Context, idx chan<- int, in <-chan S, ins []chan<- S) error {
	defer func() {
		close(idx)

		for _, i := range ins {
			close(i)
		}
	}()

	cases := make([]reflect.SelectCase, 0, len(ins)+1)

	cases = append(cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	})

	for _, c := range ins {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(c),
		})
	}

	for el := range in {
		for i := range cases[1:] {
			cases[i+1].Send = reflect.ValueOf(el)
		}

		chosen, _, _ := reflect.Select(cases)

		if chosen == 0 {
			return context.Cause(ctx)
		}

		idx <- chosen - 1 // -1 because the first case is ctx.Done()
	}

	return nil
}

func stickySpread[S any](ctx context.Context, idx chan<- int, in <-chan S, ins []chan<- S, bucket func(S) int) error {
	defer func() {
		close(idx)

		for _, i := range ins {
			close(i)
		}
	}()

	for el := range in {
		i := bucket(el)

		if err := task.Emit(ctx, ins[i], el); err != nil {
			return err
		}

		idx <- i
	}

	return nil
}
