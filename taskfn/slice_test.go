package taskfn_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/taskfn"
)

func TestSliceToSlice(t *testing.T) {
	t.Parallel()

	fn := taskfn.SliceToSlice(pipe.Map(func(ctx context.Context, el string) (string, error) {
		return el + el, nil
	}).AsTask())

	out, err := fn(context.Background(), []string{"item1", "item2"})
	require.NoError(t, err)

	assert.Equal(t, []string{"item1item1", "item2item2"}, out)
}

func TestSourceToSlice(t *testing.T) {
	t.Parallel()

	fn := taskfn.SourceToSlice(source.FromSlice([]string{"item1", "item2"}).AsSource())

	out, err := fn(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []string{"item1", "item2"}, out)
}

func TestSliceToSink(t *testing.T) {
	t.Parallel()

	out := make([]string, 0)
	fn := taskfn.SliceToSink(sink.ToSlice(&out).AsSink())

	err := fn(context.Background(), []string{"item1", "item2"})
	require.NoError(t, err)

	assert.Equal(t, []string{"item1", "item2"}, out)
}
