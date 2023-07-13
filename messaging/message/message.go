package message

import "encoding"

type Message interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}
