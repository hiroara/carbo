package pipe

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A Pipe operator that receives slices and emits those elements one by one.
type FlattenSliceOp[S any] struct {
	operator[[]S, S]
}

// Create a FlattenSlice operator.
func FlattenSlice[S any]() *FlattenSliceOp[S] {
	op := &FlattenSliceOp[S]{}
	op.operator.run = op.run
	return op
}

func (op *FlattenSliceOp[S]) run(ctx context.Context, in <-chan []S, out chan<- S) error {
	for els := range in {
		for _, el := range els {
			if err := task.Emit(ctx, out, el); err != nil {
				return err
			}
		}
	}
	return nil
}
