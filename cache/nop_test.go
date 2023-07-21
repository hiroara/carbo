package cache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/cache"
)

func TestNopSpec(t *testing.T) {
	store := cache.NewMemoryStore[string]()
	keyFn := func(s string) (string, error) {
		return "key:" + s, nil
	}
	sp := cache.NewNopSpec[string, string, string](store, keyFn)

	k, err := sp.Key("item1")
	require.NoError(t, err)
	assert.Equal(t, "key:item1", k)

	v, err := sp.Encode("item2")
	require.NoError(t, err)
	assert.Equal(t, "item2", v)

	v, err = sp.Decode("item3")
	require.NoError(t, err)
	assert.Equal(t, "item3", v)
}
