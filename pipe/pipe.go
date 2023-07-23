package pipe

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/hiroara/carbo/deferrer"
	"github.com/hiroara/carbo/task"
)

type PipeFn[S, T any] func(ctx context.Context, in <-chan S, out chan<- T) error

type Pipe[S any, T any] struct {
	deferrer.Deferrer
	run PipeFn[S, T]
}

func FromFn[S any, T any](fn PipeFn[S, T]) *Pipe[S, T] {
	return &Pipe[S, T]{run: fn}
}

func (p *Pipe[S, T]) AsTask() task.Task[S, T] {
	return task.Task[S, T](p)
}

func (p *Pipe[S, T]) Run(ctx context.Context, in <-chan S, out chan<- T) error {
	defer close(out)
	defer p.RunDeferred()
	return p.run(ctx, in, out)
}

func Concurrent[S, T any](ps []*Pipe[S, T]) *Pipe[S, T] {
	return FromFn(func(ctx context.Context, in <-chan S, out chan<- T) error {
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

func ConcurrentFromFn[S, T any](fn PipeFn[S, T], concurrency int) *Pipe[S, T] {
	ps := make([]*Pipe[S, T], concurrency)
	for i := range ps {
		ps[i] = FromFn(fn)
	}
	return Concurrent(ps)
}
