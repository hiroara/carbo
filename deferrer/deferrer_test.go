package deferrer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hiroara/carbo/deferrer"
)

func TestDeferrer(t *testing.T) {
	t.Parallel()

	d := &deferrer.Deferrer{}
	called1 := false
	called2 := false
	d.Defer(func() { called1 = true })
	d.Defer(func() { called2 = true })

	assert.False(t, called1)
	assert.False(t, called2)

	d.RunDeferred()

	assert.True(t, called1)
	assert.True(t, called2)
}
