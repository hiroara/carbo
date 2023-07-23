package task

import (
	"context"
)

// Task is an interface that represents a component of a data pipeline.
//
// Each task takes an input channel and an output channel, and communicates with other tasks through them.
//
// Type parameters:
//   S: Type of elements fed by an upstream task
//   T: Type of elements that are passed to a downstream task
type Task[S, T any] interface {
	// Run this task.
	// Inputs for this task are provided via the `in` channel,
	// and outputs of this task should be passed to a downstream task by sending them to the `out` channel.
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
