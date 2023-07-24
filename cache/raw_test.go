package cache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/cache"
	"github.com/hiroara/carbo/cache/store"
)

func TestRawSpec(t *testing.T) {
	store := store.NewMemoryStore[string]()
	keyFn := func(s string) (*cache.StoreKey[string], error) {
		return cache.Key("key:" + s), nil
	}
	sp := cache.NewRawSpec[string, string, string](store, keyFn)

	k, err := sp.Key("item1")
	require.NoError(t, err)
	assert.NotNil(t, k)

	v, err := sp.Encode("item2")
	require.NoError(t, err)
	assert.Equal(t, "item2", v)

	v, err = sp.Decode("item3")
	require.NoError(t, err)
	assert.Equal(t, "item3", v)
}
