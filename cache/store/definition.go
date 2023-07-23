package store

import "context"

type Definition[K, V any] interface {
	Get(ctx context.Context, key K) (value Value[V], err error)
	Set(ctx context.Context, key K, value V) (err error)
}
