package sink

import (
	"context"

	"github.com/hiroara/carbo/task"
)

type ToSliceOp[S any] struct {
	result *[]S
}

func ToSlice[S any](s *[]S) *ToSliceOp[S] {
	return &ToSliceOp[S]{result: s}
}

func (op *ToSliceOp[S]) AsSink() *Sink[S] {
	result := *op.result
	return ElementWise(func(ctx context.Context, s S) error {
		result = append(result, s)
		*op.result = result
		return nil
	}).AsSink()
}

func (op *ToSliceOp[S]) AsTask() task.Task[S, struct{}] {
	return op.AsSink()
}
