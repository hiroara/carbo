package marshal

type Spec[S any] interface {
	Marshal(v S) ([]byte, error)
	Unmarshal([]byte) (S, error)
}
