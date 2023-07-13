package marshal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/messaging/marshal"
)

type dummyStruct struct {
	Value string
}

func TestGob(t *testing.T) {
	data := &dummyStruct{Value: "dummy data"}
	msg := marshal.Gob(data)
	raw, ok := msg.(*marshal.GobMessage[*dummyStruct])
	assert.True(t, ok)
	assert.Equal(t, data, raw.Value)

	bs, err := raw.MarshalBinary()
	require.NoError(t, err)
	assert.NotEmpty(t, bs)

	anotherMsg := &marshal.GobMessage[*dummyStruct]{}
	err = anotherMsg.UnmarshalBinary(bs)
	require.NoError(t, err)
	assert.Equal(t, data, anotherMsg.Value)
}
