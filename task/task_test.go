package task_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/task"
)

var double = func(ctx context.Context, in <-chan string, out chan<- string) error {
	for el := range in {
		if err := task.Emit(ctx, out, el+el); err != nil {
			return err
		}
	}
	return nil
}

func TestTaskRun(t *testing.T) {
	t.Parallel()

	in := make(chan string)
	go func() {
		defer close(in)
		in <- "item1"
		in <- "item2"
	}()

	out := make(chan string, 2)
	err := task.FromFn(double).Run(context.Background(), in, out)
	require.NoError(t, err)

	assert.Equal(t, "item1item1", <-out)
	assert.Equal(t, "item2item2", <-out)

	_, ok := <-out
	assert.False(t, ok)
}

func TestTaskRunWithInputOptions(t *testing.T) {
	t.Parallel()

	in := make(chan string)

	go func() {
		defer close(in)
		time.Sleep(1 * time.Second)
		in <- "item1"
		in <- "item2"
	}()

	out := make(chan string, 2)
	err := task.FromFn(
		double,
		task.WithInputOptions(task.WithTimeout(1*time.Nanosecond)),
	).Run(context.Background(), in, out)
	require.ErrorIs(t, err, context.DeadlineExceeded)

	_, ok := <-out
	assert.False(t, ok)
}

func TestTaskRunWithOutputOptions(t *testing.T) {
	t.Parallel()

	in := make(chan string, 2)
	in <- "item1"
	in <- "item2"

	out := make(chan string)
	go func() {
		time.Sleep(1 * time.Second)
		<-out
		<-out
	}()

	err := task.FromFn(
		double,
		task.WithOutputOptions(task.WithTimeout(1*time.Nanosecond)),
	).Run(context.Background(), in, out)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}
