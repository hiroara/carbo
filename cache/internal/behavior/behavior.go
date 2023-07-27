package behavior

import "context"

type Behavior[S, T any] interface {
	Run(context.Context, S, CacheableFn[S, T]) (T, error)
}

type Entry[T any] interface {
	Get(context.Context) (*T, error)
	Set(context.Context, T) error
}

type CacheableFn[S, T any] func(context.Context, S) (T, error)

type BehaviorType int

const (
	CacheType BehaviorType = iota
	WriteOnlyType
	BypassType
)

func New[S, T any](entry Entry[T], t BehaviorType) Behavior[S, T] {
	switch t {
	case WriteOnlyType:
		return &writeOnlyBehavior[S, T]{entry: entry}
	case BypassType:
		return &bypassBehavior[S, T]{entry: entry}
	default:
		fallthrough // Fallback to CacheType
	case CacheType:
		return &cacheBehavior[S, T]{entry: entry}
	}
}
