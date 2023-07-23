package cache

import "github.com/hiroara/carbo/cache/store"

type rawSpec[S, T, K any] struct {
	store.Store[K, T]
	keyFn KeyFn[S, K]
}

func NewRawSpec[S, T, K any](cs store.Definition[K, T], keyFn KeyFn[S, K]) Spec[S, T, K, T] {
	return &rawSpec[S, T, K]{
		Store: store.Build(cs),
		keyFn: keyFn,
	}
}

func (sp *rawSpec[S, T, K]) Key(el S) (*StoreKey[K], error) {
	return sp.keyFn(el)
}

func (sp *rawSpec[S, T, K]) Encode(v T) (T, error) {
	return v, nil
}

func (sp *rawSpec[S, T, K]) Decode(v T) (T, error) {
	return v, nil
}
