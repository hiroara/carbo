package sink_test

import (
	"context"
	"testing"

	"github.com/hiroara/carbo/sink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createArraySink() (sink.SinkFn[string], *[]string, *bool) {
	items := make([]string, 0)
	called := false
	sinkFn := func(ctx context.Context, in <-chan string) error {
		called = true
		for i := range in {
			items = append(items, i)
		}
		return nil
	}
	return sinkFn, &items, &called
}

func TestSinkRun(t *testing.T) {
	t.Parallel()

	sinkFn, items, _ := createArraySink()
	sink := sink.FromFn(sinkFn)

	in := make(chan string, 2)
	out := make(chan struct{}, 1)
	in <- "item1"
	in <- "item2"
	close(in)

	err := sink.Run(context.Background(), in, out)
	require.NoError(t, err)

	assert.Equal(t, []string{"item1", "item2"}, *items)
}

func TestConcurrentSink(t *testing.T) {
	t.Parallel()

	runSink := func(s *sink.Sink[string]) error {
		in := make(chan string, 4)
		out := make(chan struct{})
		in <- "item1"
		in <- "item2"
		in <- "item3"
		in <- "item4"
		close(in)

		return s.Run(context.Background(), in, out)
	}

	t.Run("Concurrent", func(t *testing.T) {
		sinkFn1, items1, called1 := createArraySink()
		sinkFn2, items2, called2 := createArraySink()

		s := sink.Concurrent([]*sink.Sink[string]{
			sink.FromFn(sinkFn1),
			sink.FromFn(sinkFn2),
		})

		err := runSink(s)
		require.NoError(t, err)

		items := append(*items1, *items2...)
		assert.True(t, *called1)
		assert.True(t, *called2)
		assert.ElementsMatch(t, []string{"item1", "item2", "item3", "item4"}, items)
	})

	t.Run("ConcurrentFromFn", func(t *testing.T) {
		sinkFn, items, called := createArraySink()

		s := sink.ConcurrentFromFn(sinkFn, 2)

		err := runSink(s)
		require.NoError(t, err)
		assert.True(t, *called)
		assert.ElementsMatch(t, []string{"item1", "item2", "item3", "item4"}, *items)
	})
}
