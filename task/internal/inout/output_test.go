package inout_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/task/internal/inout"
)

func TestOutput(t *testing.T) {
	t.Parallel()

	dest := make(chan string)
	out := inout.NewOutput(dest, nil)
	src := out.Chan()

	go func() {
		defer close(src)
		src <- "string1"
		src <- "string2"
	}()

	checked := make(chan struct{})
	go func() {
		defer close(checked)

		assert.Equal(t, "string1", <-dest)
		assert.Equal(t, "string2", <-dest)
	}()

	require.NoError(t, inout.StartWithContext[string](context.Background(), out))
	<-checked // Wait until consumer goroutine is done
}

func TestOutputWithTimeout(t *testing.T) {
	t.Parallel()

	dest := make(chan string)
	out := inout.NewOutput(dest, &inout.Options{Timeout: 1 * time.Nanosecond})

	// Slow downstream
	go func() {
		time.Sleep(10 * time.Second)
		<-dest
	}()

	src := out.Chan()
	go func() {
		defer close(src)
		src <- "item1"
	}()

	err := inout.StartWithContext[string](context.Background(), out)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}
