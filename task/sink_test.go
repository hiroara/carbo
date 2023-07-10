package task_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/task"
)

func createArraySink() (task.SinkFn[string], *[]string, *bool) {
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
	sink := task.SinkFromFn(sinkFn)

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

	sinkFn1, items1, called1 := createArraySink()
	sinkFn2, items2, called2 := createArraySink()

	sink := task.ConcurrentSink([]*task.Sink[string]{
		task.SinkFromFn(sinkFn1),
		task.SinkFromFn(sinkFn2),
	})

	in := make(chan string, 4)
	out := make(chan struct{})
	in <- "item1"
	in <- "item2"
	in <- "item3"
	in <- "item4"
	close(in)

	err := sink.Run(context.Background(), in, out)
	require.NoError(t, err)

	items := append(*items1, *items2...)
	assert.True(t, *called1)
	assert.True(t, *called2)
	assert.ElementsMatch(t, []string{"item1", "item2", "item3", "item4"}, items)
}
