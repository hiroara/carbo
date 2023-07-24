package cache

import (
	"context"

	"github.com/hiroara/carbo/cache/internal/behavior"
)

// A signature of a function that the caching is applicable.
type CacheableFn[S, T any] func(context.Context, S) (T, error)

// Run the passed CacheableFn with caching.
//
// The function's result is cached depending on the passed Spec and its argument.
func Run[S, T, K, V any](ctx context.Context, sp Spec[S, T, K, V], arg S, fn CacheableFn[S, T]) (T, error) {
	key, err := sp.Key(arg)
	if err != nil {
		var zero T
		return zero, err
	}

	return getBehavior(sp, key).Run(ctx, arg, behavior.CacheableFn[S, T](fn))
}
