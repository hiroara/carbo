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
	"github.com/hiroara/carbo/taskfn"
)

func TestPull(t *testing.T) {
	t.Parallel()

	ms := marshal.Bytes[string]()
	data := []string{"message1", "message2", "message3"}

	// Expose data
	dir := t.TempDir()
	sock := filepath.Join(dir, "srv.sock")
	lis, err := net.Listen("unix", sock)
	require.NoError(t, err)

	expose := taskfn.SliceToSink(sink.Expose(lis, ms, 0).AsSink())

	ctx := context.Background()
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error { return expose(ctx, data) })

	// Pull data
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "unix", addr)
		}),
	}
	conn, err := grpc.Dial(sock, dialOpts...)
	require.NoError(t, err)

	out := make([]string, 0)
	pull := task.Connect(
		source.Pull(conn, ms, 0).AsTask(),
		sink.ToSlice(&out).AsTask(),
		0,
	)

	grp.Go(func() error { return flow.FromTask(pull).Run(ctx) })

	require.NoError(t, grp.Wait())

	assert.Equal(t, data, out)
}
