package source

import (
	"context"

	"github.com/hiroara/carbo/task"
)

type source[T any] struct {
	run SourceFn[T]
}

type Source[T any] interface {
	task.Task[struct{}, T]
	AsTask() task.Task[struct{}, T]
}

type SourceFn[T any] func(ctx context.Context, out chan<- T) error

func FromFn[T any](fn SourceFn[T]) Source[T] {
	return &source[T]{run: fn}
}

func (s *source[T]) AsTask() task.Task[struct{}, T] {
	return task.Task[struct{}, T](s)
}

func (s *source[T]) Run(ctx context.Context, in <-chan struct{}, out chan<- T) error {
	<-in // Initial input channel will be closed immediately after starting the flow
	defer close(out)
	return s.run(ctx, out)
}
