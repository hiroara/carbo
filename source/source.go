package source

import (
	"context"

	"github.com/hiroara/carbo/deferrer"
	"github.com/hiroara/carbo/task"
)

type source[T any] struct {
	deferrer.Deferrer
	run SourceFn[T]
}

// A task that is used as a starting point of a data pipeline.
//
// A Source usually generates multiple elements and feeds them to a downstream task,
// for example, from a provided slice or by reading a file.
//
// The input channel needs to be closed without sending any elements.
type Source[T any] task.Task[struct{}, T]

// A function that defines a Source's behavior.
// This function should send elements to the passed output channel.
// Please note that this function should not close the output channel.
type SourceFn[T any] func(ctx context.Context, out chan<- T) error

// Build a Source with a SourceFn.
func FromFn[T any](fn SourceFn[T]) Source[T] {
	return &source[T]{run: fn}
}

// Convert the Source as a task.
func (s *source[T]) AsTask() task.Task[struct{}, T] {
	return task.Task[struct{}, T](s)
}

// Run this Source.
func (s *source[T]) Run(ctx context.Context, in <-chan struct{}, out chan<- T) error {
	<-in // Initial input channel will be closed immediately after starting the flow
	defer close(out)
	defer s.RunDeferred()
	return s.run(ctx, out)
}
