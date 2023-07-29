package taskfn

import (
	"context"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
)

// Function that takes an input slice and returns an output slice.
type SliceToSliceFn[S, T any] func(context.Context, []S) ([]T, error)

// Convert a task into SliceToSliceFn.
// This is useful in case that a task takes a few length of inputs and
// returns a few length of outputs, for example, for testing a task.
func SliceToSlice[S, T any](t task.Task[S, T]) SliceToSliceFn[S, T] {
	return func(ctx context.Context, inputs []S) ([]T, error) {
		in := task.Connect(
			source.FromSlice(inputs).AsTask(),
			t,
			0,
		)

		outputs := make([]T, 0)
		out := task.Connect(
			in,
			sink.ToSlice(&outputs).AsTask(),
			0,
		)

		if err := flow.FromTask(out).Run(ctx); err != nil {
			return nil, err
		}
		return outputs, nil
	}
}
