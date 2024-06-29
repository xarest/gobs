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

	for {
		chosen, value, ok := reflect.Select(cases)
		if chosen == 0 {
			// Context was cancelled
			chErr <- ctx.Err()
			return
		} else if ok {
			// An event was received
			if err := onEvents(ctx, value.Interface().(T)); err != nil {
				return
			}
		} else {
			// This channel was closed, remove it from the slice
			cases = append(cases[:chosen], cases[chosen+1:]...)
			if len(cases) == 1 { // Only the context's Done channel is left
				return
			}
		}
	}
}
