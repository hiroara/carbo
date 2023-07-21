package source

import (
	"context"
	"reflect"

	"golang.org/x/sync/errgroup"
)

func Concurrent[T any](ss []Source[T]) Source[T] {
	return FromFn(func(ctx context.Context, out chan<- T) error {
		grp, ctx := errgroup.WithContext(ctx)
		outs := make([]chan T, len(ss))
		for i := range ss {
			s := ss[i]
			o := make(chan T)
			grp.Go(func() error {
				in := make(chan struct{})
				close(in)
				return s.Run(ctx, in, o)
			})
			outs[i] = o
		}
		grp.Go(func() error {
			cases := make([]reflect.SelectCase, len(ss))
			for i, o := range outs {
				cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(o)}
			}
			for len(cases) > 0 {
				chosen, recv, recvOK := reflect.Select(cases)
				if !recvOK {
					cases = append(cases[:chosen], cases[chosen+1:]...)
					continue
				}
				out <- recv.Interface().(T)
			}
			return nil
		})
		return grp.Wait()
	})
}
