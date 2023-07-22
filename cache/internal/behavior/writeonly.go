package behavior

import "context"

type writeOnlyBehavior[S, T any] struct {
	entry Entry[T]
}

func (b *writeOnlyBehavior[S, T]) Run(ctx context.Context, el S, fn CacheableFn[S, T]) (T, error) {
	var zero T

	v, err := fn(ctx, el)
	if err != nil {
		return zero, err
	}

	if err := b.entry.Set(ctx, v); err != nil {
		return zero, err
	}

	return v, nil
}
