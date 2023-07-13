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
	"github.com/hiroara/carbo/messaging/marshal"
	"github.com/hiroara/carbo/pb"
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

	src := source.FromSlice([]string{"item1", "item2"})
	exp := sink.Expose(lis, marshal.Raw[string], 2)
	conn := task.Connect(src.AsTask(), exp.AsTask(), 2)

	grp, ctx := errgroup.WithContext(context.Background())

	grp.Go(func() error {
		return flow.FromTask(conn).Run(ctx)
	})

	var resp *pb.BatchPullResponse

	grp.Go(func() error {
		grpcConn, err := grpc.Dial(sock, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "unix", addr)
		}))
		require.NoError(t, err)

		cli := pb.NewCommunicatorClient(grpcConn)
		resp, err = cli.BatchPull(ctx, &pb.BatchPullRequest{Limit: 3})
		require.NoError(t, err)

		return nil
	})

	require.NoError(t, grp.Wait())
	require.Len(t, resp.Messages, 2)

	var raw marshal.RawMessage[string]
	raw.UnmarshalBinary(resp.Messages[0].Value)
	assert.Equal(t, "item1", raw.Value)
	raw.UnmarshalBinary(resp.Messages[1].Value)
	assert.Equal(t, "item2", raw.Value)
}
