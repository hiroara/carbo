package sink_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/internal/testutils"
	"github.com/hiroara/carbo/sink"
)

func runSink(s sink.Sink[string]) error {
	in := make(chan string, 4)
	out := make(chan struct{})
	in <- "item1"
	in <- "item2"
	in <- "item3"
	in <- "item4"
	close(in)

	return s.Run(context.Background(), in, out)
}

func TestConcurrent(t *testing.T) {
	t.Parallel()

	t.Run("NormalCase", func(t *testing.T) {
		t.Parallel()

		sinkFn1, items1, called1 := createArraySink()
		sinkFn2, items2, called2 := createArraySink()

		s := sink.Concurrent([]sink.Sink[string]{
			sink.FromFn(sinkFn1),
			sink.FromFn(sinkFn2),
		})

		err := runSink(s)
		require.NoError(t, err)
		close(items1)
		close(called1)
		close(items2)
		close(called2)

		items := append(testutils.ReadItems(items1), testutils.ReadItems(items2)...)
		assert.Len(t, testutils.ReadItems(called1), 1)
		assert.Len(t, testutils.ReadItems(called2), 1)
		assert.ElementsMatch(t, []string{"item1", "item2", "item3", "item4"}, items)
	})

	t.Run("NoConcurrentSinksCase", func(t *testing.T) {
		t.Parallel()

		assert.PanicsWithValue(
			t,
			"at least 1 concurrent sink is required",
			func() { sink.Concurrent([]sink.Sink[string]{}) },
		)
	})
}

func TestConcurrentFromFn(t *testing.T) {
	t.Run("NormalCase", func(t *testing.T) {
		t.Parallel()

		sinkFn, items, called := createArraySink()

		s := sink.ConcurrentFromFn(sinkFn, 2)

		err := runSink(s)
		require.NoError(t, err)
		close(items)
		close(called)

		assert.Len(t, testutils.ReadItems(called), 2)
		assert.ElementsMatch(t, []string{"item1", "item2", "item3", "item4"}, testutils.ReadItems(items))
	})

	t.Run("ZeroConcurrencyCase", func(t *testing.T) {
		t.Parallel()

		fn, _, _ := createArraySink()

		assert.PanicsWithValue(
			t,
			"at least 1 concurrent sink is required",
			func() { sink.ConcurrentFromFn(fn, 0) },
		)
	})

	t.Run("NegativeConcurrencyCase", func(t *testing.T) {
		t.Parallel()

		fn, _, _ := createArraySink()

		assert.PanicsWithValue(
			t,
			"at least 1 concurrent sink is required",
			func() { sink.ConcurrentFromFn(fn, -1) },
		)
	})
}
