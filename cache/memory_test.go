package cache_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/cache"
)

func TestMemoryStore(t *testing.T) {
	store := cache.NewMemoryStore[string]()
	ctx := context.Background()
	key1 := "key-1"
	value1 := "value-1"

	_, ok, err := store.Get(ctx, key1)
	require.NoError(t, err)
	assert.False(t, ok)

	err = store.Set(ctx, key1, value1)
	require.NoError(t, err)

	v, ok, err := store.Get(ctx, key1)
	require.NoError(t, err)
	if assert.True(t, ok) {
		assert.Equal(t, value1, v)
	}
}
