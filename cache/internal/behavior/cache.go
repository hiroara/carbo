package behavior

import "context"

type cacheBehavior[S, T any] struct {
	entry Entry[T]
}

func (b *cacheBehavior[S, T]) Run(ctx context.Context, el S, fn CacheableFn[S, T]) (T, error) {
	var zero T

	v, ok, err := b.entry.Get(ctx)
	if err != nil {
		return zero, err
	}

	if ok {
		return v, nil
	}

	v, err = fn(ctx, el)
	if err != nil {
		return zero, err
	}

	if err := b.entry.Set(ctx, v); err != nil {
		return zero, err
	}

	return v, nil
}
