package cache

import (
	"github.com/hiroara/carbo/cache/store"
	"github.com/hiroara/carbo/messaging/marshal"
)

type marshalSpec[S, T, K any] struct {
	store.Store[K, []byte]
	keyFn     KeyFn[S, K]
	valueSpec marshal.Spec[T]
}

func NewMarshalSpec[S, T, K any](cs store.Definition[K, []byte], keyFn KeyFn[S, K], valueSpec marshal.Spec[T]) Spec[S, T, K, []byte] {
	return &marshalSpec[S, T, K]{
		Store:     store.Build(cs),
		keyFn:     keyFn,
		valueSpec: valueSpec,
	}
}

func (sp *marshalSpec[S, T, K]) Key(el S) (*StoreKey[K], error) {
	return sp.keyFn(el)
}

func (sp *marshalSpec[S, T, K]) Encode(v T) ([]byte, error) {
	return sp.valueSpec.Marshal(v)
}

func (sp *marshalSpec[S, T, K]) Decode(bs []byte) (T, error) {
	return sp.valueSpec.Unmarshal(bs)
}
