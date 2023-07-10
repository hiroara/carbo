package task

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Sink[S any] struct {
	run SinkFn[S]
}

type SinkFn[S any] func(ctx context.Context, in <-chan S) error

func SinkFromFn[S any](fn SinkFn[S]) *Sink[S] {
	return &Sink[S]{run: fn}
}

func (s *Sink[S]) AsTask() Task[S, struct{}] {
	return Task[S, struct{}](s)
}

func (s *Sink[S]) Run(ctx context.Context, in <-chan S, out chan<- struct{}) error {
	defer close(out)
	return s.run(ctx, in)
}

func ConcurrentSink[S any](ss []*Sink[S]) *Sink[S] {
	return SinkFromFn(func(ctx context.Context, in <-chan S) error {
		grp, ctx := errgroup.WithContext(ctx)
		for _, s := range ss {
			src := s
			grp.Go(func() error {
				return src.run(ctx, in)
			})
		}
		return grp.Wait()
	})
}
