package store

import "context"

type Store[K, V any] interface {
	Get(ctx context.Context, key K) (value Value[V], err error)
	Set(ctx context.Context, key K, value V) (err error)
}
