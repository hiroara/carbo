package task

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Source[T any] struct {
	run SourceFn[T]
}

type SourceFn[T any] func(ctx context.Context, out chan<- T) error

func SourceFromFn[T any](fn SourceFn[T]) *Source[T] {
	return &Source[T]{run: fn}
}

func (s *Source[T]) AsTask() Task[struct{}, T] {
	return Task[struct{}, T](s)
}

func (s *Source[T]) Run(ctx context.Context, in <-chan struct{}, out chan<- T) error {
	<-in // Initial input channel will be closed immediately after starting the flow
	defer close(out)
	return s.run(ctx, out)
}

func ConcurrentSource[T any](ss []*Source[T]) *Source[T] {
	return SourceFromFn(func(ctx context.Context, out chan<- T) error {
		grp, ctx := errgroup.WithContext(ctx)
		for _, s := range ss {
			src := s
			grp.Go(func() error {
				return src.run(ctx, out)
			})
		}
		return grp.Wait()
	})
}
