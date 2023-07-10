package task_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/task"
)

func createSourceFn(outputs []string) task.SourceFn[string] {
	return func(ctx context.Context, out chan<- string) error {
		for _, item := range outputs {
			out <- item
		}
		return nil
	}
}

func TestSourceRun(t *testing.T) {
	t.Parallel()

	src := task.SourceFromFn(createSourceFn([]string{"item1", "item2"}))

	in := make(chan struct{})
	out := make(chan string, 2)
	close(in)

	err := src.Run(context.Background(), in, out)
	require.NoError(t, err)

	assert.Equal(t, "item1", <-out)
	assert.Equal(t, "item2", <-out)
}

func TestConcurrentSource(t *testing.T) {
	t.Parallel()

	src := task.ConcurrentSource([]*task.Source[string]{
		task.SourceFromFn(createSourceFn([]string{"item1", "item2"})),
		task.SourceFromFn(createSourceFn([]string{"item3", "item4"})),
	})

	in := make(chan struct{})
	out := make(chan string, 4)
	close(in)

	err := src.Run(context.Background(), in, out)
	require.NoError(t, err)

	outputs := make([]string, 0)
	for item := range out {
		outputs = append(outputs, item)
	}
	assert.ElementsMatch(t, []string{"item1", "item2", "item3", "item4"}, outputs)
}
