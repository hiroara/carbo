package cache

import (
	"context"
	"sync"
)

type MemoryStore[V any] struct {
	cache sync.Map
}

func NewMemoryStore[V any]() *MemoryStore[V] {
	return &MemoryStore[V]{cache: sync.Map{}}
}

func (cs *MemoryStore[V]) Get(ctx context.Context, key string) (V, bool, error) {
	var zero V
	v, ok := cs.cache.Load(key)
	if !ok {
		return zero, false, nil
	}
	return v.(V), true, nil
}

func (cs *MemoryStore[V]) Set(ctx context.Context, key string, value V) error {
	cs.cache.Store(key, value)
	return nil
}
