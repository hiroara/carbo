package sink

import (
	"context"
)

// A Sink task that adds elements fed by its upstream task to the passed slice.
type ToSliceOp[S any] struct {
	ElementWiseOp[S]
}

// Create a ToSliceOp from a slice.
func ToSlice[S any](s *[]S) *ToSliceOp[S] {
	op := &ToSliceOp[S]{
		ElementWiseOp: *ElementWise(func(ctx context.Context, el S) error {
			*s = append(*s, el)
			return nil
		}),
	}
	return op
}
