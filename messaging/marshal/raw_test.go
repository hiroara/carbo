package marshal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/messaging/marshal"
)

func TestRaw(t *testing.T) {
	data := "dummy data"
	r := marshal.Raw[string]()
	bs, err := r.Marshal(data)
	require.NoError(t, err)
	assert.Equal(t, data, string(bs))

	d, err := r.Unmarshal(bs)
	require.NoError(t, err)
	assert.Equal(t, data, d)
}
