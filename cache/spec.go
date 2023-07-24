package cache

import (
	"context"

	"github.com/hiroara/carbo/cache/internal/behavior"
	"github.com/hiroara/carbo/cache/internal/entry"
	"github.com/hiroara/carbo/cache/store"
)

// Specification of a cache behavior.
type Spec[S, T, K, V any] interface {
	store.Store[K, V]            // Store that should be used as a cache store
	Key(S) (*StoreKey[K], error) // A function that converts an argument into a cache key.
	Encode(T) (V, error)         // A function that encodes a cacheable function's result into a value that will be stored in a cache store.
	Decode(V) (T, error)         // A function that decodes a stored value in a cache store into a cacheable function's result.
}

type spec[S, T, K, V any] struct {
	Spec[S, T, K, V]
}

func (sp *spec[S, T, K, V]) Get(ctx context.Context, key K) (*V, error) {
	return sp.Spec.Get(ctx, key)
}

func getBehavior[S, T, K, V any](sp Spec[S, T, K, V], k *StoreKey[K]) behavior.Behavior[S, T] {
	ent := entry.New[T, K, V](&spec[S, T, K, V]{sp}, k.key)

	return behavior.New[S](behavior.Entry[T](ent), k.behavior)
}
