package task

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/hiroara/carbo/deferrer"
	"github.com/hiroara/carbo/task/internal/inout"
	"github.com/hiroara/carbo/task/internal/metadata"
)

// Task is an interface that represents a component of a data pipeline.
//
// Each task takes an input channel and an output channel, and communicates with other tasks through them.
//
// Type parameters:
//
//	S: Type of elements fed by an upstream task
//	T: Type of elements that are passed to a downstream task
type Task[S, T any] interface {
	// Run this task.
	// Inputs for this task are provided via the `in` channel,
	// and outputs of this task should be passed to a downstream task by sending them to the `out` channel.
	//
	// This function must finish when the passed context has been canceled because the context will be canceled
	// when a downstream task has finished without consuming all elements in its input channel.
	// For example, sending a value to the input channel can block a goroutine when the channel buffer is full.
	// When a downstream task has finished without consuming all elements in its input channel, it is possible
	// that an upstream task still runs without knowing its downstream task is already finished, and keeps trying
	// to send values to its input channel. In this case, the upstream task can get stuck because of a full input channel.
	// Checking ctx.Done() with the `select` clause when sending a value to the output channel is a solution
	// to avoid this issue.
	// Please see Emit because the function is an easy shorthand to do this.
	Run(ctx context.Context, in <-chan S, out chan<- T) error

	// Convert this task as a task.
	// Usually, calling this function of a task returns the task itself.
	// This is used to cast other types of tasks into a task with proper type parameters.
	AsTask() Task[S, T]

	// Add a function that needs to be called after this task has completed.
	// More specifically, the registered function will be called just before the Run function returns its result.
	// This can be used, for example, to close a file or a database connection when this task has completed.
	Defer(func())
}

type task[S, T any] struct {
	deferrer.Deferrer
	TaskFn[S, T]
	*options
}

// A function that defines a Task's behavior.
// For more details, please see the Run function defined as a part of the Task interface.
// Please note that this function should close the output channel when the task finishes
// because task.FromFn does not automatically close the channel.
// The whole pipeline will be aborted when the returned error is not nil.
type TaskFn[S, T any] func(ctx context.Context, in <-chan S, out chan<- T) error

// Build a Task with a TaskFn.
func FromFn[S, T any](fn TaskFn[S, T], opts ...Option) Task[S, T] {
	tOpts := &options{}
	for _, opt := range opts {
		opt(tOpts)
	}
	return &task[S, T]{TaskFn: fn, options: tOpts}
}

func (t *task[S, T]) AsTask() Task[S, T] {
	return t
}

// GetName gets the current task name from a context.
// If the task is runnining as a part of a task, this returns the most closest task's name.
var GetName = metadata.GetName

func (t *task[S, T]) Run(ctx context.Context, in <-chan S, out chan<- T) error {
	defer t.RunDeferred()
	ctx = metadata.WithName(ctx, t.name)

	grp, ctx := errgroup.WithContext(ctx)

	ip := inout.NewInput(in, newOptions(t.inOpts))
	op := inout.NewOutput(out, newOptions(t.outOpts))

	grp.Go(func() error {
		err := inout.StartWithContext[S](ctx, ip)
		return ignoreIfErrDownstreamFinished(err)
	})

	grp.Go(func() error {
		err := t.TaskFn(ctx, ip.Chan(), op.Chan())
		return ignoreIfErrDownstreamFinished(err)
	})

	grp.Go(func() error {
		err := inout.StartWithContext[T](ctx, op)
		return ignoreIfErrDownstreamFinished(err)
	})

	return grp.Wait()
}
