package cache

import (
	"context"
)

type Store[K, V any] interface {
	Get(ctx context.Context, key K) (value V, ok bool, err error)
	Set(ctx context.Context, key K, value V) (err error)
}
