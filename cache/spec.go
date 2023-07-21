package cache

type Spec[S, T, K, V any] interface {
	Store[K, V]
	Key(S) (K, error)
	Encode(T) (V, error)
	Decode(V) (T, error)
}

type KeyFn[S, K any] func(S) (K, error)
