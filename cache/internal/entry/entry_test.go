package entry_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/cache"
	"github.com/hiroara/carbo/cache/internal/entry"
)

func buildSpec() entry.Spec[string, string, []byte] {
	store := cache.NewMemoryStore[[]byte]()
	return &dummySpec{Store: store}
}

type dummySpec struct {
	cache.Store[string, []byte]
}

func (sp *dummySpec) Decode(value []byte) (string, error) {
	return string(value), nil
}

func (sp *dummySpec) Encode(value string) ([]byte, error) {
	return []byte(value), nil
}

func TestEntry(t *testing.T) {
	t.Parallel()

	sp := buildSpec()
	ent := entry.New(sp, "key1")

	ctx := context.Background()

	v, ok, err := ent.Get(ctx)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Zero(t, v)

	err = ent.Set(ctx, "value1")
	require.NoError(t, err)

	v, ok, err = ent.Get(ctx)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "value1", v)
}
