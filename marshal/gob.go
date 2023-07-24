package marshal

import (
	"bytes"
	"encoding/gob"
)

type gobSpec[S any] struct{}

func Gob[S any]() Spec[S] {
	return &gobSpec[S]{}
}

func (g *gobSpec[S]) Marshal(v S) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *gobSpec[S]) Unmarshal(data []byte) (S, error) {
	var v S
	buf := bytes.NewBuffer(data)
	err := gob.NewDecoder(buf).Decode(&v)
	if err != nil {
		var zero S
		return zero, err
	}
	return v, nil
}
