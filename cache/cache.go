package cache

import (
	"context"

	"github.com/hiroara/carbo/cache/internal/behavior"
)

type CacheableFn[S, T any] behavior.CacheableFn[S, T]

func Run[S, T, K, V any](ctx context.Context, sp Spec[S, T, K, V], el S, fn CacheableFn[S, T]) (T, error) {
	key, err := sp.Key(el)
	if err != nil {
		var zero T
		return zero, err
	}

	return getBehavior(sp, key).Run(ctx, el, behavior.CacheableFn[S, T](fn))
}
