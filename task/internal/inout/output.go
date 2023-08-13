package inout

import "context"

type Output[T any] struct {
	src     chan T
	dest    chan<- T
	options *Options
}

func NewOutput[T any](c chan<- T, opts *Options) *Output[T] {
	if opts == nil {
		opts = &Options{}
	}
	src := make(chan T)
	return &Output[T]{src: src, dest: c, options: opts}
}

func (out *Output[T]) Chan() chan<- T {
	return out.src
}

func (out *Output[T]) Close() error {
	close(out.dest)
	return nil
}

func (out *Output[T]) passThrough(ctx context.Context) (bool, error) {
	el, ok := <-out.src
	if !ok {
		return false, nil
	}

	cancel := func() {}
	if out.options.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, out.options.Timeout)
	}
	defer cancel()
	select {
	case <-ctx.Done():
		return false, context.Cause(ctx)
	case out.dest <- el:
		return ok, nil
	}
}
