package source_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/taskfn"
)

func TestConcurrentSource(t *testing.T) {
	t.Parallel()

	src := source.Concurrent([]source.Source[string]{
		source.FromFn(createSourceFn([]string{"item1", "item2"})),
		source.FromFn(createSourceFn([]string{"item3", "item4"})),
	})

	fn := taskfn.SourceToSlice(src)

	out, err := fn(context.Background())
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"item1", "item2", "item3", "item4"}, out)
}
