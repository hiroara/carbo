package marshal

// Specification of marshaling.
type Spec[S any] interface {
	// Convert a provided value into a byte array
	Marshal(v S) ([]byte, error)

	// Convert a provided a byte array into a value.
	// This should be an inverse operation of Marshal.
	Unmarshal([]byte) (S, error)
}
