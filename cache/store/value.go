package store

// A type that represents a cache hit with a cached value or a cache miss.
type Value[V any] *V

// A Value that represents an existing cached value in a cache store.
func Hit[V any](v V) Value[V] {
	return &v
}

// A Value that represents a cache miss.
func Miss[V any]() Value[V] {
	return nil
}
