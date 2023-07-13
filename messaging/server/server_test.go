package server_test

import (
	"context"
	"net"
	"path/filepath"
	"testing"

	"github.com/hiroara/carbo/messaging/marshal"
	"github.com/hiroara/carbo/messaging/server"
	"github.com/hiroara/carbo/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestServer(t *testing.T) {
	t.Parallel()

	ms := marshal.Raw[string]()

	dir := t.TempDir()
	sock := filepath.Join(dir, "srv.sock")
	lis, err := net.Listen("unix", sock)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	grp, ctx := errgroup.WithContext(ctx)

	srv := server.New(lis, 2)
	grp.Go(func() error { return srv.Run(ctx) })

	// Input
	grp.Go(func() error {
		defer cancel()
		bs, err := ms.Marshal("message1")
		require.NoError(t, err)
		srv.Feed(bs)
		bs, err = ms.Marshal("message2")
		require.NoError(t, err)
		srv.Feed(bs)
		bs, err = ms.Marshal("message3")
		require.NoError(t, err)
		srv.Feed(bs)
		return nil
	})

	// Output
	out := make([]string, 0)
	grp.Go(func() error {
		resp, err := srv.BatchPull(context.Background(), &pb.BatchPullRequest{})
		if err != nil {
			return err
		}
		for _, msg := range resp.Messages {
			item, err := ms.Unmarshal(msg.Value)
			require.NoError(t, err)
			out = append(out, item)
		}
		return nil
	})

	require.NoError(t, grp.Wait())

	assert.Equal(t, []string{"message1", "message2", "message3"}, out)
}
