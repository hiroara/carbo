package source

import (
	"context"

	"google.golang.org/grpc"

	"github.com/hiroara/carbo/messaging/marshal"
	"github.com/hiroara/carbo/pb"
	"github.com/hiroara/carbo/task"
)

type PullOp[T any] struct {
	chunkSize   int
	client      pb.CommunicatorClient
	marshalSpec marshal.Spec[T]
}

func Pull[T any](conn grpc.ClientConnInterface, m marshal.Spec[T], chunkSize int) *PullOp[T] {
	return &PullOp[T]{chunkSize: chunkSize, client: pb.NewCommunicatorClient(conn), marshalSpec: m}
}

func (op *PullOp[T]) AsTask() task.Task[struct{}, T] {
	return task.Task[struct{}, T](FromFn(func(ctx context.Context, out chan<- T) error {
		for {
			resp, err := op.client.BatchPull(ctx, &pb.BatchPullRequest{Limit: int32(op.chunkSize)})
			if err != nil {
				return err
			}
			for _, msg := range resp.Messages {
				bs, err := op.marshalSpec.Unmarshal(msg.Value)
				if err != nil {
					return nil
				}
				out <- bs
			}
			if resp.Closed {
				return nil
			}
		}
	}))
}
