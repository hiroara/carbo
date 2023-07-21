package testutils

func ReadItems[T any](c <-chan T) []T {
	ss := make([]T, 0)
	for s := range c {
		ss = append(ss, s)
	}
	return ss
}
