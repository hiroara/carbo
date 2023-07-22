package behavior

import "context"

type bypassBehavior[S, T any] struct {
	entry Entry[T]
}

func (b *bypassBehavior[S, T]) Run(ctx context.Context, el S, fn CacheableFn[S, T]) (T, error) {
	return fn(ctx, el)
}
