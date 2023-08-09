package channel_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/hiroara/carbo/internal/channel"
)

func TestDuplicateOutChan(t *testing.T) {
	t.Parallel()

	t.Run("NormalCase", func(t *testing.T) {
		t.Parallel()

		out := make(chan string, 2)
		outs, agg := channel.DuplicateOutChan(out, 2)

		ctx := context.Background()
		grp, ctx := errgroup.WithContext(ctx)

		grp.Go(func() error {
			defer close(outs[0])
			outs[0] <- "out from first goroutine"
			return nil
		})

		grp.Go(func() error {
			defer close(outs[1])
			outs[1] <- "out from second goroutine"
			return nil
		})

		grp.Go(func() error { return agg(ctx) })

		require.NoError(t, grp.Wait())
		close(out)

		els := make([]string, 0)
		for el := range out {
			els = append(els, el)
		}

		assert.ElementsMatch(t, []string{"out from first goroutine", "out from second goroutine"}, els)
	})

	t.Run("ErrorInDownstreamCase", func(t *testing.T) {
		t.Parallel()

		out := make(chan string)
		outs, agg := channel.DuplicateOutChan(out, 2)

		ctx := context.Background()
		grp, ctx := errgroup.WithContext(ctx)

		grp.Go(func() error {
			defer close(outs[0])
			select {
			case outs[0] <- "out from first goroutine":
			case <-ctx.Done():
			}
			return nil
		})

		grp.Go(func() error {
			defer close(outs[1])
			select {
			case outs[1] <- "out from second goroutine":
			case <-ctx.Done():
			}
			return nil
		})

		aggCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel only aggregate function

		grp.Go(func() error { return agg(aggCtx) })

		require.Error(t, grp.Wait(), context.Canceled)
	})

	t.Run("ArgN=Zero", func(t *testing.T) {
		t.Parallel()

		out := make(chan string, 2)
		assert.PanicsWithValue(
			t,
			"argument n must be a positive value but received 0",
			func() { channel.DuplicateOutChan(out, 0) },
		)
	})
}
