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

	_ = inout.StartWithContext[string](context.Background(), out)

	assert.Equal(t, "string1", <-dest)
	assert.Equal(t, "string2", <-dest)
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

	ctx := context.Background()
	ctx = inout.StartWithContext[string](ctx, out)

	src := out.Chan()
	src <- "item1"
	close(src)

	select {
	case <-ctx.Done(): // Returned context is canceled when timeout is exceeded.
		assert.ErrorIs(t, ctx.Err(), context.Canceled)
	case <-time.After(time.Second):
		require.Fail(t, "Test timeout")
	}
}
