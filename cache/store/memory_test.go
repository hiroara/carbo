package store_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/cache/store"
)

func TestMemoryStore(t *testing.T) {
	t.Parallel()

	cs := store.NewMemoryStore[string]()
	ctx := context.Background()
	key1 := "key-1"
	value1 := "value-1"

	vp, err := cs.Get(ctx, key1)
	require.NoError(t, err)
	assert.Equal(t, store.Miss[string](), vp)

	err = cs.Set(ctx, key1, value1)
	require.NoError(t, err)

	vp, err = cs.Get(ctx, key1)
	require.NoError(t, err)
	assert.Equal(t, store.Hit(value1), vp)
}
