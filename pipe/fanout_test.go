package pipe_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFanout(t *testing.T) {
	src := source.FromSlice([]int{1, 2, 3})

	fo := pipe.Fanout[int](func(ctx context.Context, ss []string) ([]string, error) {
		return ss, nil
	})
	fo.Add(pipe.Map(func(ctx context.Context, i int) (string, error) {
		return strconv.FormatInt(int64(i), 10), nil
	}).AsTask(), 2, 2)
	fo.Add(pipe.Map(func(ctx context.Context, i int) (string, error) {
		return strconv.FormatInt(int64(i*10), 10), nil
	}).AsTask(), 2, 2)

	conn := task.Connect(src.AsTask(), fo.AsTask(), 2)

	out := make([][]string, 0)
	err := flow.FromTask(task.Connect(conn, sink.ToSlice(&out).AsTask(), 2)).Run(context.Background())
	require.NoError(t, err)

	assert.Equal(t, [][]string{{"1", "10"}, {"2", "20"}, {"3", "30"}}, out)
}
