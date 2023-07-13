package marshal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/messaging/marshal"
)

func TestRaw(t *testing.T) {
	data := "dummy data"
	msg := marshal.Raw(data)
	raw, ok := msg.(*marshal.RawMessage[string])
	assert.True(t, ok)
	assert.Equal(t, data, raw.Value)

	bs, err := raw.MarshalBinary()
	require.NoError(t, err)
	assert.Equal(t, data, string(bs))

	anotherRaw := &marshal.RawMessage[string]{}
	err = anotherRaw.UnmarshalBinary(bs)
	require.NoError(t, err)
	assert.Equal(t, data, anotherRaw.Value)
}
