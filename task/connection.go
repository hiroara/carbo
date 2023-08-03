package task

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

// Connection is a task that represents connected two tasks.
//
// Type parameters:
//
//	S: Type of elements fed by an upstream task
//	M: Type of elements that are sent from Src to Dest
//	T: Type of elements that are passed to a downstream task
type Connection[S, M, T any] struct {
	Src     Task[S, M] // The first task that is contained in this Connection.
	Dest    Task[M, T] // The second task that is contained in this Connection.
	srcOut  chan M
	destOut chan T
}

// Connect two tasks as a Connection.
func Connect[S, M, T any](src Task[S, M], dest Task[M, T], buf int, opts ...Option) Task[S, T] {
	conn := &Connection[S, M, T]{Src: src, Dest: dest, srcOut: make(chan M, buf), destOut: make(chan T)}
	return FromFn(conn.run, opts...)
}

var errDownstreamFinished = errors.New("a downstream task has finished")

// Run two tasks that the Connection contains.
func (conn *Connection[S, M, T]) run(ctx context.Context, in <-chan S, out chan<- T) error {
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error { return conn.Src.Run(ctx, in, conn.srcOut) })

	// destOut will be closed by Dest.
	grp.Go(func() error { return conn.Dest.Run(ctx, conn.srcOut, conn.destOut) })

	// out will be closed by *task.Run.
	grp.Go(func() error {
		for el := range conn.destOut {
			if err := Emit(ctx, out, el); err != nil {
				return err
			}
		}
		return errDownstreamFinished
	})

	err := grp.Wait()
	if errors.Is(err, errDownstreamFinished) {
		err = nil
	}
	return err
}
