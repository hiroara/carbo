package sink_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/internal/testutils"
	"github.com/hiroara/carbo/sink"
)

func createArraySink() (sink.SinkFn[string], chan string, chan struct{}) {
	c := make(chan string, 4)
	called := make(chan struct{}, 2)
	sinkFn := func(ctx context.Context, in <-chan string) error {
		called <- struct{}{}
		for i := range in {
			c <- i
		}
		return nil
	}
	return sinkFn, c, called
}

func TestSinkRun(t *testing.T) {
	t.Parallel()

	sinkFn, items, called := createArraySink()
	sink := sink.FromFn(sinkFn)

	in := make(chan string, 2)
	out := make(chan struct{}, 1)
	in <- "item1"
	in <- "item2"
	close(in)

	err := sink.Run(context.Background(), in, out)
	require.NoError(t, err)
	close(items)
	close(called)

	assert.Equal(t, []string{"item1", "item2"}, testutils.ReadItems(items))
	assert.Len(t, testutils.ReadItems(called), 1)
}
