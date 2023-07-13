package sink

import (
	"context"
	"net"

	"golang.org/x/sync/errgroup"

	"github.com/hiroara/carbo/messaging/marshal"
	"github.com/hiroara/carbo/messaging/server"
	"github.com/hiroara/carbo/task"
)

type ExposeOp[S any] struct {
	server      *server.Server
	marshalSpec marshal.Spec[S]
}

func Expose[S any](lis net.Listener, m marshal.Spec[S], buffer int) *ExposeOp[S] {
	return &ExposeOp[S]{
		server:      server.New(lis, buffer),
		marshalSpec: m,
	}
}

func (op *ExposeOp[S]) AsTask() task.Task[S, struct{}] {
	return FromFn(func(ctx context.Context, in <-chan S) error {
		grp, ctx := errgroup.WithContext(ctx)
		ctx, cancel := context.WithCancel(ctx)
		grp.Go(func() error { return op.server.Run(ctx) })
		grp.Go(func() error {
			defer cancel()
			for el := range in {
				bs, err := op.marshalSpec.Marshal(el)
				if err != nil {
					return err
				}
				op.server.Feed(bs)
			}
			return nil
		})
		return grp.Wait()
	})
}
