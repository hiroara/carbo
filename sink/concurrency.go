package sink

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/hiroara/carbo/task"
)

// Create a Sink from multiple Sinks.
// The passed Sinks will run concurrently.
func Concurrent[S any](ss []Sink[S], opts ...task.Option) Sink[S] {
	if len(ss) == 0 {
		panic("at least 1 concurrent sink is required")
	}

	return FromFn(func(ctx context.Context, in <-chan S) error {
		grp, ctx := errgroup.WithContext(ctx)
		for _, s := range ss {
			sin := s
			grp.Go(func() error {
				return sin.Run(ctx, in, make(chan<- struct{}))
			})
		}
		return grp.Wait()
	}, opts...)
}

// Create a Sink to run the provided SinkFn concurrently.
// This is a shorthand to create a concurrent Sink from Sinks with the same function.
func ConcurrentFromFn[S any](fn SinkFn[S], concurrency int, opts ...task.Option) Sink[S] {
	if concurrency < 0 {
		concurrency = 0
	}
	ss := make([]Sink[S], concurrency)
	for i := range ss {
		ss[i] = FromFn(fn, opts...)
	}
	return Concurrent(ss)
}
