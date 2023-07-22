package cache_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hiroara/carbo/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Parallel()

	createSpec := func(keyFn cache.KeyFn[string, string]) cache.Spec[string, string, string, string] {
		cs := cache.NewMemoryStore[string]()
		return cache.NewRawSpec[string, string, string](cs, keyFn)
	}

	t.Run("NormalCase", func(t *testing.T) {
		t.Parallel()

		sp := createSpec(func(el string) (*cache.StoreKey[string], error) {
			return cache.Key("key:" + el), nil
		})

		called := 0
		fn := func(ctx context.Context, el string) (string, error) {
			called += 1
			return el + el, nil
		}

		ctx := context.Background()

		v, err := cache.Run(ctx, sp, "item1", fn)
		require.NoError(t, err)
		assert.Equal(t, "item1item1", v)
		assert.Equal(t, 1, called)

		v, err = cache.Run(ctx, sp, "item1", fn)
		require.NoError(t, err)
		assert.Equal(t, "item1item1", v)
		assert.Equal(t, 1, called) // Not increased because the result is cached
	})

	t.Run("ErrorCase", func(t *testing.T) {
		t.Parallel()

		keyFnErr := errors.New("test error")
		sp := createSpec(func(el string) (*cache.StoreKey[string], error) {
			return nil, keyFnErr
		})

		fn := func(ctx context.Context, el string) (string, error) {
			return el + el, nil
		}

		ctx := context.Background()

		_, err := cache.Run(ctx, sp, "item1", fn)
		assert.ErrorIs(t, err, keyFnErr)
	})
}
