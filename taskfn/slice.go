package taskfn

import (
	"context"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
)

// Function that takes an input slice and returns an output slice or an error.
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

// Function that does not take any arguments and returns an output slice or an error.
type ToSliceFn[T any] func(context.Context) ([]T, error)

// Convert a source task into ToSliceFn.
// This is useful in case that a task does not take any inputs and
// returns a few length of outputs, for example, for testing a source task.
func SourceToSlice[T any](src source.Source[T]) ToSliceFn[T] {
	return func(ctx context.Context) ([]T, error) {
		outputs := make([]T, 0)
		out := task.Connect(
			src.AsTask(),
			sink.ToSlice(&outputs).AsTask(),
			0,
		)

		if err := flow.FromTask(out).Run(ctx); err != nil {
			return nil, err
		}
		return outputs, nil
	}
}

// Function that takes an input slice and returns nothing or an error.
type SliceFn[S any] func(context.Context, []S) error

// Convert a sink task into SliceFn.
// This is useful in case that a task takes a few length of inputs and
// returns nothing as outputs, for example, for testing a sink task.
func SliceToSink[S any](sin sink.Sink[S]) SliceFn[S] {
	return func(ctx context.Context, inputs []S) error {
		out := task.Connect(
			source.FromSlice(inputs).AsTask(),
			sin.AsTask(),
			0,
		)

		return flow.FromTask(out).Run(ctx)
	}
}
