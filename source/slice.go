package source

import (
	"context"

	"github.com/hiroara/carbo/task"
)

type SliceSourceOp[T any] struct {
	run SourceFn[T]
}

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

func (op *SliceSourceOp[T]) AsTask() task.Task[struct{}, T] {
	return task.Task[struct{}, T](FromFn(op.run))
}
