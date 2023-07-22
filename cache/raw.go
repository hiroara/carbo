package cache

type rawSpec[S, T, K any] struct {
	Store[K, T]
	keyFn KeyFn[S, K]
}

func NewRawSpec[S, T, K any](store Store[K, T], keyFn KeyFn[S, K]) Spec[S, T, K, T] {
	return &rawSpec[S, T, K]{
		Store: store,
		keyFn: keyFn,
	}
}

func (sp *rawSpec[S, T, K]) Key(el S) (K, error) {
	return sp.keyFn(el)
}

func (sp *rawSpec[S, T, K]) Encode(v T) (T, error) {
	return v, nil
}

func (sp *rawSpec[S, T, K]) Decode(v T) (T, error) {
	return v, nil
}
