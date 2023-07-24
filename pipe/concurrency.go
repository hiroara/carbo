package pipe

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/hiroara/carbo/internal/channel"
)

// Create a Pipe from multiple Pipes.
// The passed Pipes will run concurrently, and those outputs will be merged as outputs of the created Pipe.
func Concurrent[S, T any](ps []Pipe[S, T]) Pipe[S, T] {
	return FromFn(func(ctx context.Context, in <-chan S, out chan<- T) error {
		grp, ctx := errgroup.WithContext(ctx)
		outs, agg := channel.DuplicateOutChan(out, len(ps))
		for i := range ps {
			p := ps[i]
			o := outs[i]
			grp.Go(func() error {
				return p.Run(ctx, in, o)
			})
		}
		grp.Go(agg)
		return grp.Wait()
	})
}

// Create a Pipe from multiple PipeFns.
// This is a shorthand to create a concurrent Pipe from Pipes with the same function.
func ConcurrentFromFn[S, T any](fn PipeFn[S, T], concurrency int) Pipe[S, T] {
	ps := make([]Pipe[S, T], concurrency)
	for i := range ps {
		ps[i] = FromFn(fn)
	}
	return Concurrent(ps)
}
