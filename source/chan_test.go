package source_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/taskfn"
)

func TestFromChan(t *testing.T) {
	t.Parallel()

	items := []string{"item1", "item2"}

	c := make(chan string, len(items))
	for _, i := range items {
		c <- i
	}
	close(c)

	src := source.FromChan(c)
	fn := taskfn.SourceToSlice(src.AsSource())

	out, err := fn(context.Background())
	require.NoError(t, err)

	assert.Equal(t, items, out)
}
