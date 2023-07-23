package task

import (
	"context"
)

type Task[S, T any] interface {
	Run(ctx context.Context, in <-chan S, out chan<- T) error
	AsTask() Task[S, T]
	Defer(func())
}
