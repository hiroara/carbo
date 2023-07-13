package marshal

import (
	"bytes"
	"encoding/gob"

	"github.com/hiroara/carbo/messaging/message"
)

type GobMessage[S any] struct {
	value S
}

type GobMarshal struct{}

func Gob[S any](v S) message.Message[S] {
	return &GobMessage[S]{value: v}
}

func (msg *GobMessage[S]) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(msg.value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (msg *GobMessage[S]) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	err := gob.NewDecoder(buf).Decode(&msg.value)
	if err != nil {
		return err
	}
	return nil
}

func (msg *GobMessage[S]) Value() S {
	return msg.value
}
