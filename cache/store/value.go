package store

type Value[V any] *V

func Hit[V any](v V) Value[V] {
	return &v
}

func Miss[V any]() Value[V] {
	return nil
}
