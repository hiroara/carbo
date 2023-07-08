package task

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Connection[T1, T2, T3 any] struct {
	Src  Task[T1, T2]
	Dest Task[T2, T3]
	c    chan T2
}

func Connect[T1, T2, T3 any](src Task[T1, T2], dest Task[T2, T3], buf int) Task[T1, T3] {
	return &Connection[T1, T2, T3]{Src: src, Dest: dest, c: make(chan T2, buf)}
}

func (c *Connection[T1, T2, T3]) AsTask() Task[T1, T3] {
	return Task[T1, T3](c)
}

func (conn *Connection[T1, T2, T3]) Run(ctx context.Context, in <-chan T1, out chan<- T3) error {
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error { return conn.Src.Run(ctx, in, conn.c) })
	grp.Go(func() error { return conn.Dest.Run(ctx, conn.c, out) })
	return grp.Wait()
}
