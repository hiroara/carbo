package task

import (
	"context"
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
	// Checking ctx.Done() with the `select` clause when sending a value to the output channel is a solution to avoid this issue.
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
