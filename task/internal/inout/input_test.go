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

	_ = in.StartWithContext(context.Background())

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

	out := make([]string, 0)
	go func() {
		defer close(src)
		time.Sleep(1 * time.Second)
		for el := range dest {
			out = append(out, el)
		}
	}()

	ctx := context.Background()
	ctx = in.StartWithContext(ctx)

	timeout := time.After(10 * time.Second)
	for {
		select {
		case src <- "string1":
		case <-ctx.Done(): // Timeout by input option
			assert.ErrorIs(t, context.Cause(ctx), context.DeadlineExceeded)
			return
		case <-timeout:
			require.Fail(t, "Test timeout")
		}
	}
}
