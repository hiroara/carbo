package marshal

type BytesCompatible interface {
	~string | []byte
}

type rawSpec[S BytesCompatible] struct{}

func Raw[S BytesCompatible]() Spec[S] {
	return &rawSpec[S]{}
}

func (r *rawSpec[S]) Marshal(v S) ([]byte, error) {
	return []byte(v), nil
}

func (r *rawSpec[S]) Unmarshal(data []byte) (S, error) {
	return S(data), nil
}
