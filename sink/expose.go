package sink

import (
	"context"
	"net"

	"golang.org/x/sync/errgroup"

	"github.com/hiroara/carbo/internal/messaging/server"
	"github.com/hiroara/carbo/marshal"
	"github.com/hiroara/carbo/task"
)

// A Sink task that exposes elements fed by its upstream task via a gRPC service.
type ExposeOp[S any] struct {
	server      *server.Server
	marshalSpec marshal.Spec[S]
}

// Create an expose operator that runs a gRPC service.
//
// The gRPC service is bounded to the passed listener,
// and each message is encoded as defined with the passed marshal spec.
//
// Use source.Pull to receive elements exposed via this operator.
// Or, it is also possible to consume elements exposed by directly communicating with the gRPC service.
// For more details on how to communicate with the service, please see the Communicator service definition.
func Expose[S any](lis net.Listener, m marshal.Spec[S], buffer int) *ExposeOp[S] {
	return &ExposeOp[S]{
		server:      server.New(lis, buffer),
		marshalSpec: m,
	}
}

// Convert the elementwise operator as a sink.
func (op *ExposeOp[S]) AsSink(opts ...task.Option) Sink[S] {
	return FromFn(func(ctx context.Context, in <-chan S) error {
		grp, ctx := errgroup.WithContext(ctx)
		grp.Go(func() error { return op.server.Run(ctx) })
		grp.Go(func() error {
			defer op.server.Close()
			for el := range in {
				bs, err := op.marshalSpec.Marshal(el)
				if err != nil {
					return err
				}
				op.server.Feed(ctx, bs)
			}
			return nil
		})
		return grp.Wait()
	}, opts...)
}

// Convert the elementwise operator as a task.
func (op *ExposeOp[S]) AsTask() task.Task[S, struct{}] {
	return op.AsSink()
}
