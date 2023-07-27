package cache

import (
	"github.com/hiroara/carbo/cache/internal/behavior"
)

// A function that defines a corresponding key for an argument.
//
// The return value should be a StoreKey that wraps a key.
type KeyFn[S, K any] func(S) (*StoreKey[K], error)

// A KeyFn that uses the argument as a key without any conversion.
//
// The expected behavior is a normal cache behavior.
func IdentityKey[S any](el S) (*StoreKey[S], error) {
	return Key(el), nil
}

// A key for a cache store.
//
// This struct wraps a key and defines a behavior for the argument,
// like a normal cache behavior, a write-only behavior, or a bypass behavior.
type StoreKey[K any] struct {
	key      K
	behavior behavior.BehaviorType
}

// This returns a StoreKey that expects a normal cache behavior.
//
// With this key, the expected behavior is:
//
//	When a cached result exists, call a cacheable function, store the result, and return it.
//	When a cached result does not exist, return the cached value.
func Key[K any](v K) *StoreKey[K] {
	return &StoreKey[K]{key: v, behavior: behavior.CacheType}
}

// This returns a StoreKey that expects a write-only cache behavior.
//
// With this key, the expected behavior is:
//
//	When a cached result exists, call a cacheable function, store the result, and return it.
//	When a cached result does not exist, call a cacheable function, overwrite the existing cache with the result, and return it.
func WriteOnlyKey[K any](v K) *StoreKey[K] {
	return &StoreKey[K]{key: v, behavior: behavior.WriteOnlyType}
}

// This returns a StoreKey that expects a bypass cache behavior.
//
// With this key, the expected behavior is:
//
//	When a cached result exists, call a cacheable function, and return the result of the function.
//	When a cached result does not exist, call a cacheable function, and return the result of the function.
func Bypass[K any]() *StoreKey[K] {
	return &StoreKey[K]{behavior: behavior.BypassType}
}
