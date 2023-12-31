package pipe

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/hiroara/carbo/task"
)

// A Pipe task that has multiple downstream tasks, and aggregates those results.
// Each input to this operator is sent to all its downstreams and processed by them,
// and those results will be passed to this operator's aggregate function.
// This operator emits elements that the aggregate function returns.
type FanoutOp[S, I, T any] struct {
	operator[S, T]
	aggregate FanoutAggregateFn[I, T]
	tasks     []task.Task[S, I]
	inputs    []chan S
	outputs   []chan I
}

// Create a fanout operator from an aggregate function.
func Fanout[S, I, T any](aggFn FanoutAggregateFn[I, T]) *FanoutOp[S, I, T] {
	op := &FanoutOp[S, I, T]{aggregate: aggFn}
	op.operator.run = op.run
	return op
}

// Create a fanout operator from a map function.
func FanoutWithMap[S, I, T any](mapFn FanoutMapFn[I, T]) *FanoutOp[S, I, T] {
	aggFn := func(ctx context.Context, els []I, out chan<- T) error {
		o, err := mapFn(ctx, els)
		if err != nil {
			return err
		}
		return task.Emit(ctx, out, o)
	}
	op := &FanoutOp[S, I, T]{aggregate: aggFn}
	op.operator.run = op.run
	return op
}

// A function to aggregate results from downstream tasks, and send outputs to the passed output channel.
// It is ensured that all of the passed elements are created from the same input,
// and the order of results is the same as the order of registered tasks.
type FanoutAggregateFn[I, T any] func(context.Context, []I, chan<- T) error

// A function to aggregate results from downstream tasks, and return an output.
// This is a variation of FanoutAggregateFn that emits only one output.
type FanoutMapFn[I, T any] func(context.Context, []I) (T, error)

// Register a task as a downstream of the fanout operator.
func (op *FanoutOp[S, I, T]) Add(t task.Task[S, I], inBuffer, outBuffer int) {
	op.tasks = append(op.tasks, t)
	op.inputs = append(op.inputs, make(chan S, inBuffer))
	op.outputs = append(op.outputs, make(chan I, outBuffer))
}

// Run this fanout operator.
func (op *FanoutOp[S, I, T]) run(ctx context.Context, in <-chan S, out chan<- T) error {
	parentCtx := ctx

	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error { return op.feed(ctx, in) })

	var stErr error
	mutex := &sync.Mutex{}
	recordStErr := func(err error) error {
		mutex.Lock()
		defer mutex.Unlock()
		if err == nil || stErr != nil {
			return nil
		}
		if errors.Is(err, errUnmatchingLength) {
			return nil
		}
		stErr = err
		return err
	}

	for i := range op.tasks {
		ic := i
		grp.Go(func() error {
			return recordStErr(op.runTask(ctx, ic))
		})
	}

	grp.Go(func() error { return op.emit(ctx, out) })

	err := grp.Wait()

	// Check why errUnmatchingLength has been thrown
	if errors.Is(err, errUnmatchingLength) {
		if parentErr := context.Cause(parentCtx); parentErr != nil {
			// If this op has been cancelled by parent, propagate the cause
			err = parentErr
		} else if stErr != nil {
			// Check if a subtask is the cause.
			err = stErr
		}
	}

	return err
}

func (op *FanoutOp[S, I, T]) feed(ctx context.Context, in <-chan S) error {
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
		doneCase := reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()),
		}
		for len(cases) > 0 {
			chosen, _, _ := reflect.Select(append(cases, doneCase))
			if chosen == len(cases) { // ctx.Done()
				return context.Cause(ctx)
			}
			cases = append(cases[:chosen], cases[chosen+1:]...)
		}
	}
	return nil
}

var errUnmatchingLength = errors.New("unmatching length of outputs detected")

func (op *FanoutOp[S, I, T]) emit(ctx context.Context, out chan<- T) error {
	for {
		var closed bool
		els := make([]I, len(op.outputs))
		for i, o := range op.outputs {
			el, ok := <-o
			if i == 0 {
				closed = !ok
			} else if closed == ok {
				return errUnmatchingLength
			} else if i == len(op.outputs)-1 && closed {
				return nil
			}
			els[i] = el
		}
		err := op.aggregate(ctx, els, out)
		if err != nil {
			return err
		}
	}
}

func (op *FanoutOp[S, I, T]) runTask(ctx context.Context, idx int) error {
	return op.tasks[idx].Run(ctx, op.inputs[idx], op.outputs[idx])
}
