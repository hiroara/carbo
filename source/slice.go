package source

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A Source task that emits elements in the passed slice.
type SliceSourceOp[T any] struct {
	run SourceFn[T]
}

// Create a SliceSourceOp from a slice.
func FromSlice[T any](items []T) *SliceSourceOp[T] {
	return &SliceSourceOp[T]{
		run: func(ctx context.Context, out chan<- T) error {
			for _, item := range items {
				out <- item
			}
			return nil
		},
	}
}

// Convert this operation as a Source.
func (op *SliceSourceOp[T]) AsSource() Source[T] {
	return FromFn(op.run)
}

// Convert this operation as a Task.
func (op *SliceSourceOp[T]) AsTask() task.Task[struct{}, T] {
	return op.AsSource()
}
