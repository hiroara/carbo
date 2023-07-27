package flow

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A type that defines an entire data pipeline.
//
// The behavior of a Flow is defined by a task that doesn't have any input and output.
// Such a task is typically built as a pipeline that starts with a Source
// and ends with a Sink.
type Flow struct {
	task task.Task[struct{}, struct{}]
}

// Create a Flow with a task that has neither input nor output.
func FromTask(task task.Task[struct{}, struct{}]) *Flow {
	return &Flow{task: task}
}

// Run this Flow.
func (f *Flow) Run(ctx context.Context) error {
	in := make(chan struct{})
	out := make(chan struct{})
	close(in) // Kick sources
	return f.task.Run(ctx, in, out)
}
