package task

import (
	"context"

	"golang.org/x/sync/errgroup"
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

func ConcurrentPipe[S, T any](ps []*Pipe[S, T]) *Pipe[S, T] {
	return PipeFromFn(func(ctx context.Context, in <-chan S, out chan<- T) error {
		grp, ctx := errgroup.WithContext(ctx)
		for _, p := range ps {
			pipe := p
			grp.Go(func() error {
				return pipe.run(ctx, in, out)
			})
		}
		return grp.Wait()
	})
}

func ConcurrentPipeFromFn[S, T any](fn PipeFn[S, T], concurrency int) *Pipe[S, T] {
	ps := make([]*Pipe[S, T], concurrency)
	for i := range ps {
		ps[i] = PipeFromFn(fn)
	}
	return ConcurrentPipe(ps)
}
