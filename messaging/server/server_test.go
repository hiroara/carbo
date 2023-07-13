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
		srv.Feed(marshal.Raw("message1"))
		srv.Feed(marshal.Raw("message2"))
		srv.Feed(marshal.Raw("message3"))
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
			raw := &marshal.RawMessage[string]{}
			require.NoError(t, raw.UnmarshalBinary(msg.Value))
			out = append(out, raw.Value)
		}
		return nil
	})

	require.NoError(t, grp.Wait())

	assert.Equal(t, []string{"message1", "message2", "message3"}, out)
}
