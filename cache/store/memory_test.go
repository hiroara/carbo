package store_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/cache/store"
)

func TestMemoryStore(t *testing.T) {
	store := store.NewMemoryStore[string]()
	ctx := context.Background()
	key1 := "key-1"
	value1 := "value-1"

	vp, err := store.Get(ctx, key1)
	require.NoError(t, err)
	assert.Nil(t, vp)

	err = store.Set(ctx, key1, value1)
	require.NoError(t, err)

	vp, err = store.Get(ctx, key1)
	require.NoError(t, err)
	if assert.NotNil(t, vp) {
		assert.Equal(t, value1, *vp)
	}
}
