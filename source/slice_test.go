package source_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/taskfn"
)

func TestFromSlice(t *testing.T) {
	t.Parallel()

	src := source.FromSlice([]string{"item1", "item2"})
	fn := taskfn.SourceToSlice(src.AsSource())

	out, err := fn(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []string{"item1", "item2"}, out)
}
