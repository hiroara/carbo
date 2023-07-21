package sink

import (
	"context"

	"golang.org/x/sync/errgroup"
)

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

func ConcurrentFromFn[S any](fn SinkFn[S], concurrency int) Sink[S] {
	ss := make([]Sink[S], concurrency)
	for i := range ss {
		ss[i] = FromFn(fn)
	}
	return Concurrent(ss)
}
