package source_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/source"
)

func TestConcurrentSource(t *testing.T) {
	t.Parallel()

	src := source.Concurrent([]source.Source[string]{
		source.FromFn(createSourceFn([]string{"item1", "item2"})),
		source.FromFn(createSourceFn([]string{"item3", "item4"})),
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
