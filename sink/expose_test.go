package sink_test

import (
	"context"
	"net"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/internal/messaging/pb"
	"github.com/hiroara/carbo/marshal"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
)

func TestExpose(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	sock := filepath.Join(dir, "srv.sock")
	lis, err := net.Listen("unix", sock)
	require.NoError(t, err)

	data := []string{"item1", "item2"}
	ms := marshal.Bytes[string]()
	src := source.FromSlice(data)
	exp := sink.Expose(lis, ms, 2)
	conn := task.Connect(src.AsTask(), exp.AsTask(), 2)

	grp, ctx := errgroup.WithContext(context.Background())

	grp.Go(func() error {
		return flow.FromTask(conn).Run(ctx)
	})

	var fbResp *pb.FillBatchResponse
	var gbResp *pb.GetBatchResponse

	grp.Go(func() error {
		dialOpts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, "unix", addr)
			}),
		}
		grpcConn, err := grpc.Dial(sock, dialOpts...)
		require.NoError(t, err)

		cli := pb.NewCommunicatorClient(grpcConn)
		fbResp, err = cli.FillBatch(ctx, &pb.FillBatchRequest{Limit: 3})
		require.NoError(t, err)
		assert.False(t, fbResp.End)

		gbResp, err = cli.GetBatch(ctx, &pb.GetBatchRequest{})
		require.NoError(t, err)

		fbResp, err = cli.FillBatch(ctx, &pb.FillBatchRequest{Limit: 3, Token: gbResp.Token})
		require.NoError(t, err)
		assert.True(t, fbResp.End)

		return nil
	})

	require.NoError(t, grp.Wait())

	require.True(t, fbResp.End)
	require.Len(t, gbResp.Messages, 2)

	for i, msg := range gbResp.Messages {
		item, err := ms.Unmarshal(msg.Value)
		if assert.NoError(t, err) {
			assert.Equal(t, data[i], item)
		}
	}
}
