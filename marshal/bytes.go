package marshal

// Type compatible with byte array
type BytesCompatible interface {
	~string | []byte
}

type bytesSpec[S BytesCompatible] struct{}

// Create a bytes spec.
//
// This Spec simply cast a BytesCompatible type into a byte array.
func Bytes[S BytesCompatible]() Spec[S] {
	return &bytesSpec[S]{}
}

func (r *bytesSpec[S]) Marshal(v S) ([]byte, error) {
	return []byte(v), nil
}

func (r *bytesSpec[S]) Unmarshal(data []byte) (S, error) {
	return S(data), nil
}
