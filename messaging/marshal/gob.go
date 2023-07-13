package marshal

import (
	"bytes"
	"encoding/gob"

	"github.com/hiroara/carbo/messaging/message"
)

type GobMessage[S any] struct {
	Value S
}

type GobMarshal struct{}

func Gob[S any](v S) message.Message {
	return &GobMessage[S]{Value: v}
}

func (msg *GobMessage[S]) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(msg.Value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (msg *GobMessage[S]) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	err := gob.NewDecoder(buf).Decode(&msg.Value)
	if err != nil {
		return err
	}
	return nil
}
