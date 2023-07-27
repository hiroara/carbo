package task

import "context"

// Emit sends an element to the provided channel.
// In addition to simply sending a value to the channel, this function also takes care of the provided context.
// When the provided context is canceled, this function returns an error that explains why it is canceled,
// without sending any value to the channel.
func Emit[T any](ctx context.Context, out chan<- T, el T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case out <- el:
		return nil
	}
}
