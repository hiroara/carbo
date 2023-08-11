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
	Src  Task[S, M] // The first task that is contained in this Connection.
	Dest Task[M, T] // The second task that is contained in this Connection.
	c    chan M
}

// Connect two tasks as a Connection.
func Connect[S, M, T any](src Task[S, M], dest Task[M, T], buf int, opts ...Option) Task[S, T] {
	conn := &Connection[S, M, T]{Src: src, Dest: dest, c: make(chan M, buf)}
	return FromFn(conn.run, opts...)
}

var ErrAbort = errors.New("connection aborted")

// Run two tasks that the Connection contains.
func (conn *Connection[S, M, T]) run(ctx context.Context, in <-chan S, out chan<- T) error {
	grp, grpctx := errgroup.WithContext(ctx)

	grp.Go(func() error { return conn.Src.Run(grpctx, in, conn.c) })

	grp.Go(func() error { return conn.Dest.Run(ctx, conn.c, out) })

	err := grp.Wait()
	if errors.Is(err, ErrAbort) {
		err = nil
	}
	return err
}
