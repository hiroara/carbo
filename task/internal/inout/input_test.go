package inout_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/task/internal/inout"
)

func TestInput(t *testing.T) {
	t.Parallel()

	src := make(chan string)
	in := inout.NewInput(src, nil)
	dest := in.Chan()

	go func() {
		defer close(src)
		src <- "string1"
		src <- "string2"
	}()

	checked := make(chan struct{})
	go func() {
		defer close(checked)

		out := make([]string, 0)
		for el := range dest {
			out = append(out, el)
		}
		assert.Equal(t, []string{"string1", "string2"}, out)
	}()

	require.NoError(t, inout.StartWithContext[string](context.Background(), in))
	<-checked // Wait until consumer goroutine is done
}

func TestInputWithTimeout(t *testing.T) {
	t.Parallel()

	src := make(chan string)
	in := inout.NewInput(src, &inout.Options{Timeout: 1 * time.Nanosecond})
	dest := in.Chan()

	// Slow upstream
	go func() {
		defer close(src)
		time.Sleep(10 * time.Second)
		src <- "string1"
	}()

	err := inout.StartWithContext[string](context.Background(), in)
	assert.ErrorIs(t, err, context.DeadlineExceeded)

	_, ok := <-dest
	assert.False(t, ok)
}
