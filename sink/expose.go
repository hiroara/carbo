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
	server     *server.Server
	marshaller marshal.Marshaller[S]
}

func Expose[S any](lis net.Listener, m marshal.Marshaller[S], buffer int) *ExposeOp[S] {
	return &ExposeOp[S]{
		server:     server.New(lis, buffer),
		marshaller: m,
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
				op.server.Feed(op.marshaller(el))
			}
			return nil
		})
		return grp.Wait()
	})
}
