package inout

type Input[T any] struct {
	*inOut[T]
	dest <-chan T
}

func NewInput[T any](c <-chan T, opts *Options) *Input[T] {
	dest := make(chan T)
	return &Input[T]{inOut: newInOut(c, dest, opts), dest: dest}
}

func (in *Input[T]) Chan() <-chan T {
	return in.dest
}
