package task

import (
	"context"
)

type PipeFn[S, T any] func(ctx context.Context, in <-chan S, out chan<- T) error

type Pipe[S any, T any] struct {
	run PipeFn[S, T]
}

func PipeFromFn[S any, T any](fn PipeFn[S, T]) *Pipe[S, T] {
	return &Pipe[S, T]{run: fn}
}

func (p *Pipe[S, T]) AsTask() Task[S, T] {
	return Task[S, T](p)
}

func (p *Pipe[S, T]) Run(ctx context.Context, in <-chan S, out chan<- T) error {
	defer close(out)
	return p.run(ctx, in, out)
}
