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

func (cs *MemoryStore[V]) Get(ctx context.Context, key string) (*V, error) {
	v, ok := cs.cache.Load(key)
	if !ok {
		return nil, nil
	}
	val := v.(V)
	return &val, nil
}

func (cs *MemoryStore[V]) Set(ctx context.Context, key string, value V) error {
	cs.cache.Store(key, value)
	return nil
}
