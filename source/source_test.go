package source_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
	"github.com/hiroara/carbo/taskfn"
)

func createSourceFn(outputs []string) source.SourceFn[string] {
	return func(ctx context.Context, out chan<- string) error {
		for _, item := range outputs {
			if err := task.Emit(ctx, out, item); err != nil {
				return err
			}
		}
		return nil
	}
}

func TestSourceRun(t *testing.T) {
	t.Parallel()

	src := source.FromFn(createSourceFn([]string{"item1", "item2"}))

	deferredCalled := false
	src.Defer(func() { deferredCalled = true })

	fn := taskfn.SourceToSlice(src)

	out, err := fn(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []string{"item1", "item2"}, out)

	assert.True(t, deferredCalled)
}
