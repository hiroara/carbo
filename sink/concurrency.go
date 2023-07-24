package sink

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Create a Sink from multiple Sinks.
// The passed Sinks will run concurrently.
func Concurrent[S any](ss []Sink[S]) Sink[S] {
	return FromFn(func(ctx context.Context, in <-chan S) error {
		grp, ctx := errgroup.WithContext(ctx)
		for _, s := range ss {
			sin := s
			grp.Go(func() error {
				return sin.Run(ctx, in, make(chan<- struct{}))
			})
		}
		return grp.Wait()
	})
}

// Create a Sink to run the provided SinkFn concurrently.
// This is a shorthand to create a concurrent Sink from Sinks with the same function.
func ConcurrentFromFn[S any](fn SinkFn[S], concurrency int) Sink[S] {
	ss := make([]Sink[S], concurrency)
	for i := range ss {
		ss[i] = FromFn(fn)
	}
	return Concurrent(ss)
}
