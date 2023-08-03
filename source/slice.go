package source

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A Source task that emits elements in the passed slice.
type SliceSourceOp[T any] struct {
	operator[T]
	items []T
}

// Create a SliceSourceOp from a slice.
func FromSlice[T any](items []T) *SliceSourceOp[T] {
	op := &SliceSourceOp[T]{items: items}
	op.operator.run = op.run
	return op
}

func (op *SliceSourceOp[T]) run(ctx context.Context, out chan<- T) error {
	for _, item := range op.items {
		if err := task.Emit(ctx, out, item); err != nil {
			return err
		}
	}
	return nil
}
