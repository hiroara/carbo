package behavior_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hiroara/carbo/cache/internal/behavior"
)

func TestNew(t *testing.T) {
	t.Parallel()

	bc := behavior.New[int, *int](nil, behavior.CacheType)
	bw := behavior.New[int, *int](nil, behavior.WriteOnlyType)
	bb := behavior.New[int, *int](nil, behavior.BypassType)
	bu := behavior.New[int, *int](nil, -1) // Unknown type

	assert.NotNil(t, bc)
	assert.NotNil(t, bw)
	assert.NotNil(t, bb)
	assert.Equal(t, bc, bu)
}
