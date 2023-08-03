package task_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
)

func TestConnectionRun(t *testing.T) {
	t.Parallel()

	src := source.FromSlice([]string{"item1", "item2"})

	outputs := make([]string, 0)
	sink := sink.ToSlice(&outputs)

	conn := task.Connect(src.AsTask(), sink.AsTask(), 0)

	deferredCalled := false
	conn.Defer(func() { deferredCalled = true })

	in := make(chan struct{})
	out := make(chan struct{})
	close(in)

	err := conn.Run(context.Background(), in, out)
	require.NoError(t, err)

	assert.Equal(t, []string{"item1", "item2"}, outputs)

	assert.True(t, deferredCalled)
}
