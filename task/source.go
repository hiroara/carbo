package task

import (
	"context"
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
