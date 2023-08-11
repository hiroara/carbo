package pipe

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A task that is used as an intermediate process of a data pipeline.
//
// A Pipe usually receives elements from an upstream task via an input channel,
// process them, and feeds them to a downstream task.
type Pipe[S, T any] task.Task[S, T]

// A function that defines a Pipe's behavior.
// This function should receive elements from the passed input channel, process them,
// and pass the results to the passed output channel.
// Please note that this function should not close the passed channels
// because pipe.FromFn automatically closes the output channel
// and closing the input channel is the upstream task's responsibility.
// The whole pipeline will be aborted when the returned error is not nil.
type PipeFn[S, T any] func(ctx context.Context, in <-chan S, out chan<- T) error

// Build a Pipe with a PipeFn.
func FromFn[S any, T any](fn PipeFn[S, T], opts ...task.Option) Pipe[S, T] {
	return task.FromFn(task.TaskFn[S, T](func(ctx context.Context, in <-chan S, out chan<- T) error {
		defer close(out)
		return fn(ctx, in, out)
	}), opts...)
}
