package inout

import "context"

type Input[T any] struct {
	src     <-chan T
	dest    chan T
	options *Options
}

func NewInput[T any](c <-chan T, opts *Options) *Input[T] {
	if opts == nil {
		opts = &Options{}
	}
	dest := make(chan T)
	return &Input[T]{src: c, dest: dest, options: opts}
}

func (in *Input[T]) Chan() <-chan T {
	return in.dest
}

func (in *Input[T]) Close() error {
	close(in.dest)
	return nil
}

func (in *Input[T]) passThrough(ctx context.Context) (bool, error) {
	cancel := func() {}
	if in.options.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, in.options.Timeout)
	}
	defer cancel()
	select {
	case <-ctx.Done():
		return false, context.Cause(ctx)
	case el, ok := <-in.src:
		if ok {
			select {
			case <-ctx.Done():
			case in.dest <- el:
			}
		}
		return ok, nil
	}
}
