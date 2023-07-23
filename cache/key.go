package cache

import (
	"github.com/hiroara/carbo/cache/internal/behavior"
)

type KeyFn[S, K any] func(S) (*StoreKey[K], error)

func IdentityKey[S any](el S) (*StoreKey[S], error) {
	return Key(el), nil
}

type StoreKey[K any] struct {
	key      K
	behavior behavior.BehaviorType
}

func Key[K any](v K) *StoreKey[K] {
	return &StoreKey[K]{key: v, behavior: behavior.CacheType}
}

func WriteOnlyKey[K any](v K) *StoreKey[K] {
	return &StoreKey[K]{key: v, behavior: behavior.WriteOnlyType}
}

func Bypass[K any]() *StoreKey[K] {
	return &StoreKey[K]{behavior: behavior.BypassType}
}
