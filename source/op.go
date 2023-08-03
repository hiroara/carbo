package source

import "github.com/hiroara/carbo/task"

type operator[T any] struct {
	run SourceFn[T]
}

// Convert this operation as a Source.
func (op *operator[T]) AsSource(opts ...task.Option) Source[T] {
	return FromFn(op.run, opts...)
}

// Convert this operation as a Task.
func (op *operator[T]) AsTask(opts ...task.Option) task.Task[struct{}, T] {
	return op.AsSource(opts...)
}
