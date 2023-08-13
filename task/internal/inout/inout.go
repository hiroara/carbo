package inout

import (
	"context"
	"io"
	"time"
)

type InOut[T any] interface {
	io.Closer
	passThrough(ctx context.Context) (bool, error)
}

type Options struct {
	Timeout time.Duration
}

func StartWithContext[T any](ctx context.Context, io InOut[T]) context.Context {
	ctx, cancel := context.WithCancelCause(ctx)
	go func() {
		defer io.Close()
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
