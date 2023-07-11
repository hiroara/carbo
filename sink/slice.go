package sink

import "github.com/hiroara/carbo/task"

type ToSliceOp[S any] struct {
	result *[]S
}

func ToSlice[S any](s *[]S) *ToSliceOp[S] {
	return &ToSliceOp[S]{result: s}
}

func (op *ToSliceOp[S]) AsTask() task.Task[S, struct{}] {
	result := *op.result
	return ElementWise(func(s S) error {
		result = append(result, s)
		*op.result = result
		return nil
	}).AsTask()
}
