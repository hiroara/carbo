package source_test

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
	"github.com/hiroara/carbo/marshal"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
)

func TestPull(t *testing.T) {
	t.Parallel()

	ms := marshal.Bytes[string]()

	dir := t.TempDir()
	sock := filepath.Join(dir, "srv.sock")
	lis, err := net.Listen("unix", sock)
	require.NoError(t, err)

	data := []string{"message1", "message2", "message3"}
	src := source.FromSlice(data)
	sin := sink.Expose(lis, ms, 3)
	exposeFlow := flow.FromTask(task.Connect(src.AsTask(), sin.AsTask(), 2))

	ctx := context.Background()
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error { return exposeFlow.Run(ctx) })

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "unix", addr)
		}),
	}
	conn, err := grpc.Dial(sock, dialOpts...)
	require.NoError(t, err)
	pull := source.Pull(conn, ms, 3)
	out := make([]string, 0)
	pullFlow := flow.FromTask(task.Connect(pull.AsTask(), sink.ToSlice(&out).AsTask(), 2))

	grp.Go(func() error { return pullFlow.Run(ctx) })

	require.NoError(t, grp.Wait())

	assert.Equal(t, data, out)
}
