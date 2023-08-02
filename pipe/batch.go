package pipe

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A Pipe task that makes a fixed size of batches.
type BatchOp[S any] struct {
	size int
}

// Create a batch operator with the passed size.
func Batch[S any](size int) *BatchOp[S] {
	return &BatchOp[S]{size: size}
}

// Convert the batch operator as a Pipe.
func (op *BatchOp[S]) AsPipe(opts ...task.Option) Pipe[S, []S] {
	return FromFn(op.run, opts...)
}

// Convert the batch operator as a task.
func (op *BatchOp[S]) AsTask() task.Task[S, []S] {
	return op.AsPipe().AsTask()
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
