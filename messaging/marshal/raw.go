package marshal

import (
	"github.com/hiroara/carbo/messaging/message"
)

type BytesCompatible interface {
	~string | []byte
}

type RawMessage[S BytesCompatible] struct {
	value S
}

func Raw[S BytesCompatible](v S) message.Message[S] {
	return &RawMessage[S]{value: v}
}

func (msg *RawMessage[S]) MarshalBinary() ([]byte, error) {
	return []byte(msg.value), nil
}

func (msg *RawMessage[S]) UnmarshalBinary(data []byte) error {
	msg.value = S(data)
	return nil
}

func (msg *RawMessage[S]) Value() S {
	return msg.value
}
