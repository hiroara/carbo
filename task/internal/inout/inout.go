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

func StartWithContext[T any](ctx context.Context, io InOut[T]) error {
	defer io.Close()
	ok := true
	var err error
	for ok {
		ok, err = io.passThrough(ctx)
	}
	return err
}
