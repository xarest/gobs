package utils

import (
	"context"
	"reflect"
)

func WaitOnEvents[T any](ctx context.Context, onEvents func(ctx context.Context, event T) error, chErr chan error, channels ...chan T) {
	if len(channels) == 0 {
		return
	}
	var cases []reflect.SelectCase
	cases = append(cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	})

	for _, ch := range channels {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}

	for len(cases) > 1 {
		chosen, value, ok := reflect.Select(cases)
		if chosen == 0 {
			// Context was cancelled
			if chErr != nil {
				chErr <- ctx.Err()
			}
			return
		} else if ok {
			// An event was received
			if err := onEvents(ctx, value.Interface().(T)); err != nil {
				return
			}
		} else {
			// This channel was closed, remove it from the slice
			cases = append(cases[:chosen], cases[chosen+1:]...)
		}
	}
}

func WrapChannel[T, K any](ctx context.Context, ch chan T) chan K {
	wrappedCh := make(chan K)
	go func() {
		defer close(wrappedCh)
		k := new(K)
		select {
		case <-ch:
		case <-ctx.Done():
		}
		wrappedCh <- *k
	}()
	return wrappedCh
}
