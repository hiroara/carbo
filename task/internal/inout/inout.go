package inout

import (
	"context"
	"time"
)

type inOut[T any] struct {
	*Options
	src  <-chan T
	dest chan<- T
}

type Options struct {
	Timeout time.Duration
}

func newInOut[T any](src <-chan T, dest chan<- T, opts *Options) *inOut[T] {
	if opts == nil {
		opts = &Options{}
	}
	return &inOut[T]{src: src, dest: dest, Options: opts}
}

func (io *inOut[T]) StartWithContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancelCause(ctx)
	go func() {
		defer close(io.dest)
		ok := true
		var err error
		for ok {
			ok, err = io.passThrough(ctx)
		}
		if err != nil {
			cancel(err)
		}
	}()
	return ctx
}

func (io *inOut[T]) passThrough(ctx context.Context) (bool, error) {
	cancel := func() {}
	if io.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, io.Timeout)
	}
	defer cancel()
	select {
	case <-ctx.Done():
		return false, context.Cause(ctx)
	case el, ok := <-io.src:
		if ok {
			if err := io.emit(ctx, el); err != nil {
				return false, err
			}
		}
		return ok, nil
	}
}

func (io *inOut[T]) emit(ctx context.Context, el T) error {
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case io.dest <- el:
		return nil
	}
}
