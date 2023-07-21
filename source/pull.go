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

func (op *PullOp[T]) AsSource() *Source[T] {
	return FromFn(func(ctx context.Context, out chan<- T) error {
		lim := int32(op.chunkSize)
		fbResp, err := op.client.FillBatch(ctx, &pb.FillBatchRequest{Limit: lim})
		if err != nil {
			return err
		}

		for !fbResp.End {
			resp, err := op.client.GetBatch(ctx, &pb.GetBatchRequest{})
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

			fbResp, err = op.client.FillBatch(ctx, &pb.FillBatchRequest{Token: resp.Token, Limit: lim})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (op *PullOp[T]) AsTask() task.Task[struct{}, T] {
	return op.AsSource()
}
