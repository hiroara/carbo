package pipe

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A Pipe operator that makes a fixed size of batches.
type BatchOp[S any] struct {
	operator[S, []S]
	size int
}

// Create a Batch operator.
func Batch[S any](size int) *BatchOp[S] {
	op := &BatchOp[S]{size: size}
	op.operator.run = op.run
	return op
}

func (op *BatchOp[S]) run(ctx context.Context, in <-chan S, out chan<- []S) error {
	b := make([]S, 0, op.size)
	for el := range in {
		b = append(b, el)
		if len(b) < op.size {
			continue
		}
		if err := task.Emit(ctx, out, b); err != nil {
			return err
		}
		b = make([]S, 0, op.size)
	}
	if len(b) > 0 {
		if err := task.Emit(ctx, out, b); err != nil {
			return err
		}
	}
	return nil
}
