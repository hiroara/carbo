package pipe_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/taskfn"
)

func TestSelect(t *testing.T) {
	t.Parallel()

	selectEven := taskfn.SliceToSlice(pipe.Select(func(n int) bool { return n%2 == 0 }).AsTask())

	result, err := selectEven(context.Background(), []int{1, 2, 3, 4, 5})
	require.NoError(t, err)
	assert.Equal(t, []int{2, 4}, result)
}
