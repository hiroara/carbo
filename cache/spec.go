package cache

import (
	"github.com/hiroara/carbo/cache/internal/behavior"
	"github.com/hiroara/carbo/cache/internal/entry"
	"github.com/hiroara/carbo/cache/store"
)

type Spec[S, T, K, V any] interface {
	store.Store[K, V]
	Key(S) (*StoreKey[K], error)
	Encode(T) (V, error)
	Decode(V) (T, error)
}

func getBehavior[S, T, K, V any](sp Spec[S, T, K, V], k *StoreKey[K]) behavior.Behavior[S, T] {
	ent := entry.New[T, K, V](sp, k.key)

	return behavior.New[S](behavior.Entry[T](ent), k.behavior)
}
