package store

import (
	"context"
	"sync"
)

// A cache store that stores cached values in memory.
//
// Thanks to sync.Map, this cache store is safe for concurrent use by multiple goroutines.
type MemoryStore[V any] struct {
	cache sync.Map
}

// Create a MemoryStore.
func NewMemoryStore[V any]() *MemoryStore[V] {
	return &MemoryStore[V]{cache: sync.Map{}}
}

// Get a cached value with a key.
func (cs *MemoryStore[V]) Get(ctx context.Context, key string) (Value[V], error) {
	v, ok := cs.cache.Load(key)
	if !ok {
		return nil, nil
	}
	val := v.(V)
	return Value[V](&val), nil
}

// Store a value with a key.
func (cs *MemoryStore[V]) Set(ctx context.Context, key string, value V) error {
	cs.cache.Store(key, value)
	return nil
}
