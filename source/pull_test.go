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

	"github.com/hiroara/carbo/marshal"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/taskfn"
)

func TestPull(t *testing.T) {
	t.Parallel()

	ms := marshal.Bytes[string]()
	data := []string{"message1", "message2", "message3"}

	dial := func(sock string) (*grpc.ClientConn, error) {
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, "unix", addr)
			}),
		}
		return grpc.Dial(sock, opts...)
	}

	t.Run("NormalCase", func(t *testing.T) {
		t.Parallel()

		// Expose data
		sock := filepath.Join(t.TempDir(), "srv.sock")
		lis, err := net.Listen("unix", sock)
		require.NoError(t, err)
		expose := taskfn.SliceToSink(sink.Expose(lis, ms, 0).AsSink())

		ctx := context.Background()
		grp, ctx := errgroup.WithContext(ctx)

		grp.Go(func() error { return expose(ctx, data) })

		// Pull data
		conn, err := dial(sock)
		require.NoError(t, err)

		pull := taskfn.SourceToSlice[string](source.Pull(conn, ms, 0).AsTask())

		out, err := pull(ctx)
		require.NoError(t, err)
		require.NoError(t, grp.Wait())

		assert.Equal(t, data, out)
	})

	t.Run("ErrorCase", func(t *testing.T) {
		t.Parallel()

		sock := filepath.Join(t.TempDir(), "srv.sock")
		lis, err := net.Listen("unix", sock)
		require.NoError(t, err)
		expose := taskfn.SliceToSink(sink.Expose(lis, ms, 0).AsSink())

		ctx := context.Background()
		grp, ctx := errgroup.WithContext(ctx)

		grp.Go(func() error { return expose(ctx, data) })

		// Pull data
		conn, err := dial(sock)
		require.NoError(t, err)

		// Pass a wrong marshal spec to cause an error
		pull := taskfn.SourceToSlice[string](source.Pull(conn, marshal.Gob[string](), 0).AsTask())

		_, err = pull(ctx)
		require.Error(t, err)
		require.NoError(t, grp.Wait())
	})
}
