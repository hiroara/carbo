package cache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/cache"
	"github.com/hiroara/carbo/messaging/marshal"
)

func TestMarshalSpec(t *testing.T) {
	store := cache.NewMemoryStore[[]byte]()
	keyFn := func(s string) (string, error) {
		return "key:" + s, nil
	}
	sp := cache.NewMarshalSpec[string, string, string](store, keyFn, marshal.Bytes[string]())

	k, err := sp.Key("item1")
	require.NoError(t, err)
	assert.Equal(t, "key:item1", k)

	bs, err := sp.Encode("item2")
	require.NoError(t, err)
	assert.Equal(t, "item2", string(bs))

	v, err := sp.Decode([]byte("item3"))
	require.NoError(t, err)
	assert.Equal(t, "item3", v)
}
