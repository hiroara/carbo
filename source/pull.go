package source

import (
	"context"

	"google.golang.org/grpc"

	"github.com/hiroara/carbo/internal/messaging/pb"
	"github.com/hiroara/carbo/marshal"
	"github.com/hiroara/carbo/task"
)

// A Source task that reads elements from a gRPC Communicator service and emits them.
//
// This Source task can be used to pull data from another process that exposes data via a Communicator service.
type PullOp[T any] struct {
	batchSize   int
	client      pb.CommunicatorClient
	marshalSpec marshal.Spec[T]
}

// Create a PullOp with a gRPC connection and a marshal spec.
//
// The gRPC connection needs to be a connection with a Communicator service.
// And, the marshal spec defines how elements exposed via the Communicator service should be decoded.
// To successfully decode elements pulled from the Communicator service, the marshal spec needs to be
// the same one used to expose the data.
//
// The batchSize argument defines the size of batches when pulling data from a Communicator service.
// The larger batch size reduces the number of communication over a network, but also,
// it let the process to wait for a large batch is fulfilled.
func Pull[T any](conn grpc.ClientConnInterface, m marshal.Spec[T], batchSize int) *PullOp[T] {
	return &PullOp[T]{batchSize: batchSize, client: pb.NewCommunicatorClient(conn), marshalSpec: m}
}

// Convert this operation as a Source.
func (op *PullOp[T]) AsSource() Source[T] {
	return FromFn(op.handleError(op.run))
}

// Convert this operation as a Task.
func (op *PullOp[T]) AsTask() task.Task[struct{}, T] {
	return op.AsSource()
}

func (op *PullOp[T]) run(ctx context.Context, out chan<- T) error {
	lim := int32(op.batchSize)
	fbResp, err := op.client.FillBatch(ctx, &pb.FillBatchRequest{Limit: lim})
	if err != nil {
		return err
	}

	for !fbResp.End {
		resp, err := op.client.GetBatch(ctx, &pb.GetBatchRequest{})
		if err != nil {
			return err
		}

		fbResp, err = op.client.FillBatch(ctx, &pb.FillBatchRequest{Token: resp.Token, Limit: lim})
		if err != nil {
			return err
		}

		for _, msg := range resp.Messages {
			el, err := op.marshalSpec.Unmarshal(msg.Value)
			if err != nil {
				return err
			}
			if err := task.Emit(ctx, out, el); err != nil {
				return err
			}
		}
	}
	return nil
}

func (op *PullOp[T]) handleError(fn SourceFn[T]) SourceFn[T] {
	return func(ctx context.Context, out chan<- T) error {
		err := fn(ctx, out)
		if err != nil {
			// Ignore error
			op.client.Abort(ctx, &pb.AbortRequest{Reason: &pb.AbortReason{Message: err.Error()}})
		}
		return err
	}
}
