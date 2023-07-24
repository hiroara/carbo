package marshal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/marshal"
)

type dummyStruct struct {
	Value string
}

func TestGob(t *testing.T) {
	data := &dummyStruct{Value: "dummy data"}
	m := marshal.Gob[*dummyStruct]()
	bs, err := m.Marshal(data)
	require.NoError(t, err)
	assert.NotEmpty(t, bs)

	d, err := m.Unmarshal(bs)
	require.NoError(t, err)
	assert.Equal(t, data, d)
}
