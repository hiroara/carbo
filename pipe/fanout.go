package pipe

import (
	"context"
	"errors"
	"reflect"

	"golang.org/x/sync/errgroup"

	"github.com/hiroara/carbo/deferrer"
	"github.com/hiroara/carbo/task"
)

// A Pipe task that has multiple downstream tasks, and aggregates those results.
// Each input to this operator is sent to all its downstreams and processed by them,
// and those results will be passed to this operator's aggregate function.
// This operator emits elements that the aggregate function returns.
type FanoutOp[S, I, T any] struct {
	deferrer.Deferrer
	aggregate FanoutAggregateFn[I, T]
	tasks     []task.Task[S, I]
	inputs    []chan S
	outputs   []chan I
}

// Create a fanout operator from an aggregate function.
func Fanout[S, I, T any](aggFn FanoutAggregateFn[I, T]) *FanoutOp[S, I, T] {
	return &FanoutOp[S, I, T]{aggregate: aggFn}
}

// A function to aggregate results from downstream tasks.
// It is ensured that all of the passed elements are created from the same input,
// and the order of results is the same as the order of registered tasks.
type FanoutAggregateFn[S, T any] func(context.Context, []S) (T, error)

// Register a task as a downstream of the fanout operator.
func (op *FanoutOp[S, I, T]) Add(t task.Task[S, I], inBuffer, outBuffer int) {
	op.tasks = append(op.tasks, t)
	op.inputs = append(op.inputs, make(chan S, inBuffer))
	op.outputs = append(op.outputs, make(chan I, outBuffer))
}

// Run this fanout operator.
func (op *FanoutOp[S, I, T]) Run(ctx context.Context, in <-chan S, out chan<- T) error {
	defer op.RunDeferred()
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error { return op.feed(in) })
	for i := range op.tasks {
		ic := i
		grp.Go(func() error {
			return op.runTask(ctx, ic)
		})
	}
	grp.Go(func() error { return op.emit(ctx, out) })
	return grp.Wait()
}

// Convert the fanout operator as a task.
func (op *FanoutOp[S, I, T]) AsTask() task.Task[S, T] {
	return task.Task[S, T](op)
}

func (op *FanoutOp[S, I, T]) feed(in <-chan S) error {
	defer func() {
		for _, ic := range op.inputs {
			close(ic)
		}
	}()
	for el := range in {
		cases := make([]reflect.SelectCase, len(op.inputs))
		for i, ic := range op.inputs {
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectSend,
				Chan: reflect.ValueOf(ic),
				Send: reflect.ValueOf(el),
			}
		}
		for len(cases) > 0 {
			chosen, _, _ := reflect.Select(cases)
			cases = append(cases[:chosen], cases[chosen+1:]...)
		}
	}
	return nil
}

func (op *FanoutOp[S, I, T]) emit(ctx context.Context, out chan<- T) error {
	defer close(out)
	for {
		var closed bool
		els := make([]I, len(op.outputs))
		for i, o := range op.outputs {
			el, ok := <-o
			if i == 0 {
				closed = !ok
			} else if closed == ok {
				return errors.New("Unmatching length of outputs detected")
			} else if i == len(op.outputs)-1 && closed {
				return nil
			}
			els[i] = el
		}
		r, err := op.aggregate(ctx, els)
		if err != nil {
			return err
		}
		out <- r
	}
}

func (op *FanoutOp[S, I, T]) runTask(ctx context.Context, idx int) error {
	return op.tasks[idx].Run(ctx, op.inputs[idx], op.outputs[idx])
}
