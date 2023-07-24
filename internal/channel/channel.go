package channel

import "reflect"

func DuplicateOutChan[T any](out chan<- T, n int) ([]chan<- T, func() error) {
	outs := make([]chan<- T, n)
	cases := make([]reflect.SelectCase, n)
	for i := range outs {
		o := make(chan T)
		outs[i] = o
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(o)}
	}
	return outs, func() error {
		for len(cases) > 0 {
			chosen, recv, recvOK := reflect.Select(cases)
			if !recvOK {
				cases = append(cases[:chosen], cases[chosen+1:]...)
				continue
			}
			out <- recv.Interface().(T)
		}
		return nil
	}
}
