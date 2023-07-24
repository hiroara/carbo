package pipe

import (
	"context"

	"github.com/hiroara/carbo/deferrer"
	"github.com/hiroara/carbo/task"
)

// A task that is used as an intermediate process of a data pipeline.
//
// A Pipe usually receives elements from an upstream task via an input channel,
// process them, and feeds them to a downstream task.
type Pipe[S, T any] task.Task[S, T]

type pipe[S any, T any] struct {
	deferrer.Deferrer
	run PipeFn[S, T]
}

// A function that defines a Pipe's behavior.
// This function should receive elements from the passed input channel, process them,
// and pass the results to the passed output channel.
// Please note that this function should not close the passed channels.
// The whole pipeline will be aborted when the returned error is not nil.
type PipeFn[S, T any] func(ctx context.Context, in <-chan S, out chan<- T) error

// Build a Pipe with a PipeFn.
func FromFn[S any, T any](fn PipeFn[S, T]) Pipe[S, T] {
	return &pipe[S, T]{run: fn}
}

// Convert the Pipe as a task.
func (p *pipe[S, T]) AsTask() task.Task[S, T] {
	return task.Task[S, T](p)
}

// Run this Pipe.
func (p *pipe[S, T]) Run(ctx context.Context, in <-chan S, out chan<- T) error {
	defer close(out)
	defer p.RunDeferred()
	return p.run(ctx, in, out)
}
