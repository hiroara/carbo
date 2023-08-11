package pipe

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A Pipe operator that emits only N elements from its upstream task.
type TakeOp[S any] struct {
	operator[S, S]
	n int
}

// Create a Take operator.
func Take[S any](n int) *TakeOp[S] {
	op := &TakeOp[S]{n: n}
	op.operator.run = op.run
	return op
}

func (op *TakeOp[S]) run(ctx context.Context, in <-chan S, out chan<- S) error {
	c := 0
	for el := range in {
		if err := task.Emit(ctx, out, el); err != nil {
			return err
		}
		c += 1
		if c == op.n {
			return task.ErrAbort
		}
	}
	return nil
}
