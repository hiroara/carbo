package pipe_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/taskfn"
)

func TestFlattenSlice(t *testing.T) {
	t.Parallel()

	flatten := taskfn.SliceToSlice(pipe.FlattenSlice[string]().AsTask())

	result, err := flatten(context.Background(), [][]string{{"item1", "item2"}, {"item3"}})
	require.NoError(t, err)
	assert.Equal(t, []string{"item1", "item2", "item3"}, result)
}
