package store

import "context"

// An interface that defines the behavior of a key-value cache store.
type Store[K, V any] interface {
	// Get a cached value with a key.
	// The returned value needs to be wraped with Value that represents hit or miss.
	Get(ctx context.Context, key K) (value Value[V], err error)

	// Store a value with a key.
	Set(ctx context.Context, key K, value V) (err error)
}
