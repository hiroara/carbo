package task

import "context"

func Feed[T any](ctx context.Context, out chan<- T, el T) error {
	select {
	case out <- el:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
