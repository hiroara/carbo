package source

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A Source task that emits elements received from the passed chan.
type ChanSourceOp[T any] struct {
	operator[T]
	c <-chan T
}

// Create a ChanSourceOp from a slice.
func FromChan[T any](c <-chan T) *ChanSourceOp[T] {
	op := &ChanSourceOp[T]{c: c}
	op.operator.run = op.run
	return op
}

func (op *ChanSourceOp[T]) run(ctx context.Context, out chan<- T) error {
	for el := range op.c {
		if err := task.Emit(ctx, out, el); err != nil {
			return err
		}
	}
	return nil
}
