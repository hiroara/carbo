package source

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/hiroara/carbo/internal/channel"
	"github.com/hiroara/carbo/task"
)

// Create a Source from multiple Sources.
// The passed Sources will run concurrently, and those outputs will be merged as outputs of the created Source.
func Concurrent[T any](ss []Source[T], opts ...task.Option) Source[T] {
	return FromFn(func(ctx context.Context, out chan<- T) error {
		grp, ctx := errgroup.WithContext(ctx)
		outs, agg := channel.DuplicateOutChan(out, len(ss))
		for i := range ss {
			s := ss[i]
			o := outs[i]
			grp.Go(func() error {
				in := make(chan struct{})
				close(in)
				return s.Run(ctx, in, o)
			})
		}
		grp.Go(func() error { return agg(ctx) })
		return grp.Wait()
	}, opts...)
}
