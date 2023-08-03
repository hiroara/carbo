package pipe

import (
	"github.com/hiroara/carbo/task"
)

// A Pipe operator struct that can be converted into a Pipe task.
type operator[S, T any] struct {
	run PipeFn[S, T]
}

// Convert the operator into a Pipe task.
func (op *operator[S, T]) AsPipe(opts ...task.Option) Pipe[S, T] {
	return FromFn(op.run, opts...)
}

// Convert the operator into a task.
func (op *operator[S, T]) AsTask(opts ...task.Option) task.Task[S, T] {
	return op.AsPipe(opts...)
}
