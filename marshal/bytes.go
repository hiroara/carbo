package marshal

type BytesCompatible interface {
	~string | []byte
}

type bytesSpec[S BytesCompatible] struct{}

func Bytes[S BytesCompatible]() Spec[S] {
	return &bytesSpec[S]{}
}

func (r *bytesSpec[S]) Marshal(v S) ([]byte, error) {
	return []byte(v), nil
}

func (r *bytesSpec[S]) Unmarshal(data []byte) (S, error) {
	return S(data), nil
}
