package store

import (
	"context"
)

type Store[K, V any] interface {
	Get(ctx context.Context, key K) (value *V, err error)
	Set(ctx context.Context, key K, value V) (err error)
}

type store[K, V any] struct {
	Definition[K, V]
}

func (cs *store[K, V]) Get(ctx context.Context, key K) (*V, error) {
	return cs.Definition.Get(ctx, key)
}

func Build[K, V any](cs Definition[K, V]) Store[K, V] {
	return &store[K, V]{cs}
}
