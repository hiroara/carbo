package message

import "encoding"

type Message[T any] interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	Value() T
}
