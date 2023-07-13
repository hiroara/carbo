package marshal

import (
	"github.com/hiroara/carbo/messaging/message"
)

type BytesCompatible interface {
	~string | []byte
}

type RawMessage[S BytesCompatible] struct {
	Value S
}

func Raw[S BytesCompatible](v S) message.Message {
	return &RawMessage[S]{Value: v}
}

func (msg *RawMessage[S]) MarshalBinary() ([]byte, error) {
	return []byte(msg.Value), nil
}

func (msg *RawMessage[S]) UnmarshalBinary(data []byte) error {
	msg.Value = S(data)
	return nil
}
