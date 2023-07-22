package entry

import "context"

type Spec[T, K, V any] interface {
	Get(context.Context, K) (V, bool, error)
	Set(context.Context, K, V) error
	Decode(V) (T, error)
	Encode(T) (V, error)
}

type Entry[T, K, V any] struct {
	spec Spec[T, K, V]
	key  K
}

func New[T, K, V any](sp Spec[T, K, V], key K) *Entry[T, K, V] {
	return &Entry[T, K, V]{spec: sp, key: key}
}

func (ent *Entry[T, K, V]) Get(ctx context.Context) (T, bool, error) {
	var zero T

	e, ok, err := ent.spec.Get(ctx, ent.key)
	if err != nil {
		return zero, false, err
	}

	if !ok {
		return zero, false, nil
	}

	v, err := ent.spec.Decode(e)
	if err != nil {
		return zero, false, err
	}

	return v, true, nil
}

func (ent *Entry[T, K, V]) Set(ctx context.Context, v T) error {
	bs, err := ent.spec.Encode(v)
	if err != nil {
		return err
	}

	return ent.spec.Set(ctx, ent.key, bs)
}
