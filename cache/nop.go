package cache

type nopSpec[S, T, K any] struct {
	Store[K, T]
	keyFn KeyFn[S, K]
}

func NewNopSpec[S, T, K any](store Store[K, T], keyFn KeyFn[S, K]) Spec[S, T, K, T] {
	return &nopSpec[S, T, K]{
		Store: store,
		keyFn: keyFn,
	}
}

func (sp *nopSpec[S, T, K]) Key(el S) (K, error) {
	return sp.keyFn(el)
}

func (sp *nopSpec[S, T, K]) Encode(v T) (T, error) {
	return v, nil
}

func (sp *nopSpec[S, T, K]) Decode(v T) (T, error) {
	return v, nil
}
