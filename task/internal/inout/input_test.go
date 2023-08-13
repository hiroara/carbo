package inout_test

import (
	"context"
	"fmt"
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

	_ = inout.StartWithContext[string](context.Background(), in)

	go func() {
		defer close(src)
		src <- "string1"
		src <- "string2"
	}()

	out := make([]string, 0)
	for el := range dest {
		out = append(out, el)
	}
	assert.Equal(t, []string{"string1", "string2"}, out)
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

	ctx := context.Background()
	ctx = inout.StartWithContext[string](ctx, in)

	for {
		select {
		case el := <-dest:
			require.Fail(t, fmt.Sprintf("Test timeout (received %s)", el))
		case <-ctx.Done(): // Timeout by input option
			assert.ErrorIs(t, context.Cause(ctx), context.DeadlineExceeded)
			return
		case <-time.After(1 * time.Second):
			require.Fail(t, "Test timeout")
		}
	}
}
