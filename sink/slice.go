package sink

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A Sink task that adds elements fed by its upstream task to the passed slice.
type ToSliceOp[S any] struct {
	result *[]S
}

// Create a ToSliceOp from a slice.
func ToSlice[S any](s *[]S) *ToSliceOp[S] {
	return &ToSliceOp[S]{result: s}
}

// Convert this operation as a Sink.
func (op *ToSliceOp[S]) AsSink(opts ...task.Option) Sink[S] {
	result := *op.result
	return ElementWise(func(ctx context.Context, s S) error {
		result = append(result, s)
		*op.result = result
		return nil
	}).AsSink(opts...)
}

// Convert this operation as a Task.
func (op *ToSliceOp[S]) AsTask() task.Task[S, struct{}] {
	return op.AsSink()
}
