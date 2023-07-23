package sink

import (
	"context"

	"github.com/hiroara/carbo/deferrer"
	"github.com/hiroara/carbo/task"
)

type sink[S any] struct {
	deferrer.Deferrer
	run SinkFn[S]
}

type Sink[S any] task.Task[S, struct{}]

type SinkFn[S any] func(ctx context.Context, in <-chan S) error

func FromFn[S any](fn SinkFn[S]) Sink[S] {
	return &sink[S]{run: fn}
}

func (s *sink[S]) AsTask() task.Task[S, struct{}] {
	return task.Task[S, struct{}](s)
}

func (s *sink[S]) Run(ctx context.Context, in <-chan S, out chan<- struct{}) error {
	defer close(out)
	defer s.RunDeferred()
	return s.run(ctx, in)
}
