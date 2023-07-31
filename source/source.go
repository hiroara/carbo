package source

import (
	"context"

	"github.com/hiroara/carbo/task"
)

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
// The whole pipeline will be aborted when the returned error is not nil.
type SourceFn[T any] func(ctx context.Context, out chan<- T) error

// Build a Source with a SourceFn.
func FromFn[T any](fn SourceFn[T], opts ...task.Option) Source[T] {
	return task.FromFn(func(ctx context.Context, in <-chan struct{}, out chan<- T) error {
		<-in // Initial input channel will be closed immediately after starting the flow
		return fn(ctx, out)
	}, opts...)
}
