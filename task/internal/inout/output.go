package inout

type Output[T any] struct {
	*inOut[T]
	src chan<- T
}

func NewOutput[T any](c chan<- T, opts *Options) *Output[T] {
	src := make(chan T)
	return &Output[T]{inOut: newInOut(src, c, opts), src: src}
}

func (in *Output[T]) Chan() chan<- T {
	return in.src
}
