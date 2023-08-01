package channel

import (
	"context"
	"reflect"

	"github.com/hiroara/carbo/task"
)

func DuplicateOutChan[T any](out chan<- T, n int) ([]chan<- T, func(context.Context) error) {
	outs := make([]chan<- T, n)
	cases := make([]reflect.SelectCase, n)
	for i := range outs {
		o := make(chan T)
		outs[i] = o
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(o)}
	}
	return outs, func(ctx context.Context) error {
		for len(cases) > 0 {
			chosen, recv, recvOK := reflect.Select(cases)
			if !recvOK {
				cases = append(cases[:chosen], cases[chosen+1:]...)
				continue
			}
			if err := task.Emit(ctx, out, recv.Interface().(T)); err != nil {
				return err
			}
		}
		return nil
	}
}
