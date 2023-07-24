package channel_test

import (
	"context"
	"testing"

	"github.com/hiroara/carbo/internal/channel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestDuplicateOutChan(t *testing.T) {
	t.Parallel()

	out := make(chan string, 2)
	outs, agg := channel.DuplicateOutChan(out, 2)
	grp, _ := errgroup.WithContext(context.Background())

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

	grp.Go(agg)

	require.NoError(t, grp.Wait())
	close(out)

	els := make([]string, 0)
	for el := range out {
		els = append(els, el)
	}

	assert.ElementsMatch(t, []string{"out from first goroutine", "out from second goroutine"}, els)
}
