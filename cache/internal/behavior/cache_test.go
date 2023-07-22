package behavior_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/cache/internal/behavior"
)

func TestCacheBehavior(t *testing.T) {
	t.Parallel()

	t.Run("NormalCase", func(t *testing.T) {
		t.Parallel()

		called := 0
		fn := func(ctx context.Context, el string) (string, error) {
			called += 1
			return el + el, nil
		}

		ent, b := createBehavior(behavior.CacheType)

		ctx := context.Background()

		v, err := b.Run(ctx, "item1", fn)
		require.NoError(t, err)
		assert.Equal(t, "item1item1", v)
		assert.Equal(t, 1, called)

		v, ok, err := ent.Get(ctx)
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "item1item1", v)

		v, err = b.Run(ctx, "item1", fn)
		require.NoError(t, err)
		assert.Equal(t, "item1item1", v)
		assert.Equal(t, 1, called) // Not called because the first result is reused
	})

	t.Run("ErrorCase", func(t *testing.T) {
		t.Parallel()

		_, b := createBehavior(behavior.CacheType)

		fnErr := errors.New("test error")
		_, err := b.Run(context.Background(), "item1", func(ctx context.Context, el string) (string, error) {
			return "", fnErr
		})

		assert.ErrorIs(t, err, fnErr)
	})
}
