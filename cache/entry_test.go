package cache_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/cache"
	"github.com/hiroara/carbo/messaging/marshal"
)

func TestEntry(t *testing.T) {
	store := cache.NewMemoryStore[[]byte]()
	keyFn := func(s string) (string, error) {
		return "key:" + s, nil
	}
	sp := cache.NewMarshalSpec[string, string, string](store, keyFn, marshal.Raw[string]())

	ent, err := cache.GetEntry(sp, "item1")
	require.NoError(t, err)

	called := 0
	runEntry := func() (string, error) {
		return ent.Run(context.Background(), func(ctx context.Context, el string) (string, error) {
			called += 1
			return "result:" + el, nil
		})
	}

	v, err := runEntry()
	require.NoError(t, err)
	assert.Equal(t, "result:item1", v)
	assert.Equal(t, 1, called)

	v, err = runEntry()
	require.NoError(t, err)
	assert.Equal(t, "result:item1", v)
	assert.Equal(t, 1, called) // Not increased because the calculation is cached
}
