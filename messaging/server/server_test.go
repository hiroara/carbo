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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer(t *testing.T) {
	t.Parallel()

	buildServer := func() *server.Server {
		dir := t.TempDir()
		sock := filepath.Join(dir, "srv.sock")
		lis, err := net.Listen("unix", sock)
		require.NoError(t, err)

		return server.New(lis, 2)
	}

	ms := marshal.Raw[string]()

	feedMessage := func(ctx context.Context, srv *server.Server, msg string) error {
		bs, err := ms.Marshal(msg)
		if err != nil {
			return err
		}
		srv.Feed(ctx, bs)
		return nil
	}

	t.Run("Normal", func(t *testing.T) {
		t.Parallel()

		srv := buildServer()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		grp, ctx := errgroup.WithContext(ctx)

		grp.Go(func() error { return srv.Run(ctx) })

		// Input
		grp.Go(func() error {
			defer srv.Close()
			require.NoError(t, feedMessage(ctx, srv, "message1"))
			require.NoError(t, feedMessage(ctx, srv, "message2"))
			require.NoError(t, feedMessage(ctx, srv, "message3"))
			return nil
		})

		// Output
		out := make([]string, 0)
		grp.Go(func() error {
			fbResp, err := srv.FillBatch(ctx, &pb.FillBatchRequest{Limit: 2})
			if err != nil {
				return err
			}
			require.False(t, fbResp.End)

			gbResp, err := srv.GetBatch(ctx, &pb.GetBatchRequest{})
			if err != nil {
				return err
			}
			assert.Len(t, gbResp.Messages, 2)
			for _, msg := range gbResp.Messages {
				item, err := ms.Unmarshal(msg.Value)
				require.NoError(t, err)
				out = append(out, item)
			}

			token := gbResp.Token

			fbResp, err = srv.FillBatch(ctx, &pb.FillBatchRequest{Token: token, Limit: 2})
			if err != nil {
				return err
			}
			require.False(t, fbResp.End)

			gbResp, err = srv.GetBatch(ctx, &pb.GetBatchRequest{})
			if err != nil {
				return err
			}
			assert.Len(t, gbResp.Messages, 1)
			for _, msg := range gbResp.Messages {
				item, err := ms.Unmarshal(msg.Value)
				require.NoError(t, err)
				out = append(out, item)
			}

			token = gbResp.Token

			fbResp, err = srv.FillBatch(ctx, &pb.FillBatchRequest{Token: token, Limit: 2})
			if err != nil {
				return err
			}
			require.True(t, fbResp.End)

			return nil
		})

		require.NoError(t, grp.Wait())

		assert.Equal(t, []string{"message1", "message2", "message3"}, out)
	})

	t.Run("TokenUnmatch", func(t *testing.T) {
		t.Parallel()

		srv := buildServer()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		grp, ctx := errgroup.WithContext(ctx)

		grp.Go(func() error { return srv.Run(ctx) })

		// Input
		grp.Go(func() error {
			defer srv.Close()
			require.NoError(t, feedMessage(ctx, srv, "message1"))
			require.NoError(t, feedMessage(ctx, srv, "message2"))
			require.NoError(t, feedMessage(ctx, srv, "message3"))
			return nil
		})

		// Output
		grp.Go(func() error {
			_, err := srv.FillBatch(ctx, &pb.FillBatchRequest{Token: "unknown", Limit: 2})
			return err
		})

		err := grp.Wait()
		require.Error(t, err)

		s, ok := status.FromError(err)
		if assert.True(t, ok) {
			assert.Equal(t, codes.InvalidArgument, s.Code())
		}
	})

	t.Run("RepeatingGet", func(t *testing.T) {
		t.Parallel()

		srv := buildServer()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		grp, ctx := errgroup.WithContext(ctx)

		grp.Go(func() error { return srv.Run(ctx) })

		// Input
		grp.Go(func() error {
			defer srv.Close()
			require.NoError(t, feedMessage(ctx, srv, "message1"))
			require.NoError(t, feedMessage(ctx, srv, "message2"))
			require.NoError(t, feedMessage(ctx, srv, "message3"))
			return nil
		})

		// Output
		grp.Go(func() error {
			defer cancel() // Call cancel to abort other goroutines

			fbResp, err := srv.FillBatch(ctx, &pb.FillBatchRequest{Limit: 2})
			if err != nil {
				return err
			}
			require.False(t, fbResp.End)

			gbResp, err := srv.GetBatch(ctx, &pb.GetBatchRequest{})
			if err != nil {
				return err
			}
			assert.Len(t, gbResp.Messages, 2)
			prevMsgs := gbResp.Messages

			// Call GetBatch again.
			gbResp, err = srv.GetBatch(ctx, &pb.GetBatchRequest{})
			if err != nil {
				return err
			}

			// Can get the same batch again.
			// This behavior is for letting a client retry in exceptional cases.
			if assert.Len(t, gbResp.Messages, len(prevMsgs)) {
				assert.Equal(t, prevMsgs[0], gbResp.Messages[0])
				assert.Equal(t, prevMsgs[1], gbResp.Messages[1])
			}

			return nil
		})

		require.NoError(t, grp.Wait())
	})
}
