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

// A task that is used as a terminal point of a data pipeline.
//
// A Sink receives elements fed by an upstream task, and consumes them,
// for example, print them into STDOUT or writing them into a file.
//
// The output channel needs to be closed without sending any elements.
type Sink[S any] task.Task[S, struct{}]

// A function that defines a Sink's behavior.
// This function should receive elements via the passed input channel.
// The whole pipeline will be aborted when the returned error is not nil.
type SinkFn[S any] func(ctx context.Context, in <-chan S) error

// Build a Sink with a SinkFn.
func FromFn[S any](fn SinkFn[S]) Sink[S] {
	return &sink[S]{run: fn}
}

// Convert the Sink as a task.
func (s *sink[S]) AsTask() task.Task[S, struct{}] {
	return task.Task[S, struct{}](s)
}

// Run this Sink.
func (s *sink[S]) Run(ctx context.Context, in <-chan S, out chan<- struct{}) error {
	defer close(out)
	defer s.RunDeferred()
	return s.run(ctx, in)
}
