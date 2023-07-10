package flow

import (
	"context"

	"github.com/hiroara/carbo/task"
)

type Flow struct {
	task task.Task[struct{}, struct{}]
}

func FromTask(task task.Task[struct{}, struct{}]) *Flow {
	return &Flow{task: task}
}

func (f *Flow) Run(ctx context.Context) error {
	in := make(chan struct{}, 0)
	out := make(chan struct{}, 0)
	close(in) // Kick sources
	return f.task.Run(ctx, in, out)
}
