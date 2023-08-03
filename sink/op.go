package sink

import "github.com/hiroara/carbo/task"

type operator[S any] struct {
	run SinkFn[S]
}

// Convert this operation as a Sink.
func (op *operator[S]) AsSink(opts ...task.Option) Sink[S] {
	return FromFn(op.run, opts...)
}

// Convert this operation as a Task.
func (op *operator[S]) AsTask(opts ...task.Option) task.Task[S, struct{}] {
	return op.AsSink(opts...)
}
