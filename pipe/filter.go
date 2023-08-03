package pipe

import (
	"context"

	"github.com/hiroara/carbo/task"
)

// A Pipe operator that emits only elements that the passed predicate function returns true.
type SelectOp[S any] struct {
	operator[S, S]
	predicate func(S) bool
}

// Create a Select operator.
func Select[S any](predicate func(S) bool) *SelectOp[S] {
	op := &SelectOp[S]{predicate: predicate}
	op.operator.run = op.run
	return op
}

func (op *SelectOp[S]) run(ctx context.Context, in <-chan S, out chan<- S) error {
	for el := range in {
		if op.predicate(el) {
			if err := task.Emit(ctx, out, el); err != nil {
				return err
			}
		}
	}
	return nil
}
