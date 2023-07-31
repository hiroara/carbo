package inout_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

	_ = out.StartWithContext(context.Background())

	assert.Equal(t, "string1", <-dest)
	assert.Equal(t, "string2", <-dest)
}

func TestOutputWithTimeout(t *testing.T) {
	t.Parallel()

	dest := make(chan string)
	out := inout.NewOutput(dest, &inout.Options{Timeout: 1 * time.Nanosecond})
	src := out.Chan()

	go func() {
		defer close(src)
		time.Sleep(1 * time.Second)
		src <- "string1"
	}()

	ctx := context.Background()
	ctx = out.StartWithContext(ctx)

	_, ok := <-dest
	assert.False(t, ok)

	<-ctx.Done() // Returned context is canceled when timeout is exceeded.
}
