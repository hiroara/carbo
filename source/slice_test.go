package source_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/source"
)

func TestFromSlice(t *testing.T) {
	t.Parallel()

	src := source.FromSlice([]string{"item1", "item2"})

	in := make(chan struct{})
	out := make(chan string, 2)
	close(in)

	err := src.AsTask().Run(context.Background(), in, out)
	require.NoError(t, err)

	assert.Equal(t, "item1", <-out)
	assert.Equal(t, "item2", <-out)
}
