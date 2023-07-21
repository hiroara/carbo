package cache

import "context"

type Entry[S, T, K, V any] struct {
	spec Spec[S, T, K, V]
	el   S
	key  K
}

type CacheableFn[S, T any] func(context.Context, S) (T, error)

func GetEntry[S, T, K, V any](sp Spec[S, T, K, V], el S) (*Entry[S, T, K, V], error) {
	key, err := sp.Key(el)
	if err != nil {
		return nil, err
	}
	return &Entry[S, T, K, V]{spec: sp, el: el, key: key}, nil
}

func (ent *Entry[S, T, K, V]) Run(ctx context.Context, fn CacheableFn[S, T]) (T, error) {
	var zero T

	v, ok, err := ent.get(ctx)
	if err != nil {
		return zero, err
	}

	if ok {
		return v, nil
	}

	v, err = fn(ctx, ent.el)
	if err != nil {
		return zero, err
	}

	if err := ent.set(ctx, v); err != nil {
		return zero, err
	}

	return v, nil
}

func (ent *Entry[S, T, K, V]) get(ctx context.Context) (T, bool, error) {
	var zero T

	bs, ok, err := ent.spec.Get(ctx, ent.key)
	if err != nil {
		return zero, false, err
	}

	if !ok {
		return zero, false, nil
	}

	v, err := ent.spec.Decode(bs)
	if err != nil {
		return zero, false, err
	}

	return v, true, nil
}

func (ent *Entry[S, T, K, V]) set(ctx context.Context, v T) error {
	bs, err := ent.spec.Encode(v)
	if err != nil {
		return err
	}

	return ent.spec.Set(ctx, ent.key, bs)
}
