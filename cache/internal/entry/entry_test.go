package entry_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/cache/internal/entry"
	"github.com/hiroara/carbo/cache/store"
)

func buildSpec() entry.Spec[string, string, []byte] {
	cs := store.NewMemoryStore[[]byte]()
	return &dummySpec{Store: store.Build[string, []byte](cs)}
}

type dummySpec struct {
	store.Store[string, []byte]
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

	v, err := ent.Get(ctx)
	require.NoError(t, err)
	assert.Nil(t, v)

	err = ent.Set(ctx, "value1")
	require.NoError(t, err)

	v, err = ent.Get(ctx)
	require.NoError(t, err)
	if assert.NotNil(t, v) {
		assert.Equal(t, "value1", *v)
	}
}
