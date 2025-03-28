package pipe

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/hiroara/carbo/internal/channel"
	"github.com/hiroara/carbo/task"
)

type concurrency[S, T any] struct {
	run PipeFn[S, T]
}

type contextKey int

const concurrencyIndex contextKey = iota

func ConcurrencyIndex(ctx context.Context) int {
	if v, ok := ctx.Value(concurrencyIndex).(int); ok {
		return v
	}
	return -1
}

// Create a concurrent Pipe from multiple operators that have the same behavior.
//
// Each input is processed whenever it is possible.
// For this reason, the concurrent Pipe doesn't preserve the order.
func (op *concurrency[S, T]) Concurrent(concurrency int, opts ...task.Option) Pipe[S, T] {
	return ConcurrentFromFn(op.run, concurrency, opts...)
}

// Create a Pipe from multiple Pipes.
// The passed Pipes will run concurrently, and those outputs will be merged as outputs of the created Pipe.
//
// Each input is processed whenever it is possible.
// For this reason, the concurrent Pipe doesn't preserve the order.
func Concurrent[S, T any](ps []Pipe[S, T], opts ...task.Option) Pipe[S, T] {
	if len(ps) == 0 {
		panic("at least 1 concurrent pipe is required")
	}

	return FromFn(func(ctx context.Context, in <-chan S, out chan<- T) error {
		grp, ctx := errgroup.WithContext(ctx)
		outs, agg := channel.DuplicateOutChan(out, len(ps))
		for i := range ps {
			p := ps[i]
			o := outs[i]

			ctx := context.WithValue(ctx, concurrencyIndex, i)

			grp.Go(func() error { return p.Run(ctx, in, o) })
		}
		grp.Go(func() error { return agg(ctx) })
		return grp.Wait()
	}, opts...)
}

// Create a Pipe to run the provided PipeFn concurrently.
// This is a shorthand to create a concurrent Pipe from Pipes with the same function.
func ConcurrentFromFn[S, T any](fn PipeFn[S, T], concurrency int, opts ...task.Option) Pipe[S, T] {
	if concurrency < 0 {
		concurrency = 0
	}
	ps := make([]Pipe[S, T], concurrency)
	for i := range ps {
		ps[i] = FromFn(fn, opts...)
	}
	return Concurrent(ps)
}
