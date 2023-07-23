package behavior

import "context"

type cacheBehavior[S, T any] struct {
	entry Entry[T]
}

func (b *cacheBehavior[S, T]) Run(ctx context.Context, el S, fn CacheableFn[S, T]) (T, error) {
	var zero T

	v, err := b.entry.Get(ctx)
	if err != nil {
		return zero, err
	}

	if v != nil {
		return *v, nil
	}

	t, err := fn(ctx, el)
	if err != nil {
		return zero, err
	}

	if err := b.entry.Set(ctx, t); err != nil {
		return zero, err
	}

	return t, nil
}
