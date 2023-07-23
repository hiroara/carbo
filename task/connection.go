package task

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/hiroara/carbo/deferrer"
)

// Connection is a task that represents connected two tasks.
//
// Type parameters:
//   S: Type of elements fed by an upstream task
//   M: Type of elements that are sent from Src to Dest
//   T: Type of elements that are passed to a downstream task
type Connection[S, M, T any] struct {
	deferrer.Deferrer
	Src  Task[S, M] // The first task that is contained in this Connection.
	Dest Task[M, T] // The second task that is contained in this Connection.
	c    chan M
}

// Connect two tasks as a Connection.
func Connect[S, M, T any](src Task[S, M], dest Task[M, T], buf int) Task[S, T] {
	return &Connection[S, M, T]{Src: src, Dest: dest, c: make(chan M, buf)}
}

// Convert the Connection as a task.
func (c *Connection[S, M, T]) AsTask() Task[S, T] {
	return Task[S, T](c)
}

// Run two tasks that the Connection contains.
func (conn *Connection[S, M, T]) Run(ctx context.Context, in <-chan S, out chan<- T) error {
	defer conn.RunDeferred()
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error { return conn.Src.Run(ctx, in, conn.c) })
	grp.Go(func() error { return conn.Dest.Run(ctx, conn.c, out) })
	return grp.Wait()
}
