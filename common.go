package gobs

import (
	"context"
	"reflect"
)

func concurrenceProcesses[T any](ctx context.Context, processes []T, proc func(ctx context.Context, p T) error, limitProcs int) error {
	if len(processes) == 0 {
		return nil
	}
	c, cancel := context.WithCancel(ctx)
	ch := make(chan error, len(processes))
	if limitProcs <= 0 {
		limitProcs = len(processes)
	}
	sem := make(chan struct{}, limitProcs)

	defer func() {
		cancel()
		close(ch)
		close(sem)
	}()
	for _, p := range processes {
		sem <- struct{}{}
		go func(c context.Context, p T) {
			defer func() {
				<-sem
			}()
			select {
			case <-c.Done():
				ch <- c.Err()
			default:
				ch <- proc(c, p)
			}
		}(c, p)
	}
	errs := make([]error, 0, len(processes))
	for range processes {
		err := <-ch
		if err != nil {
			cancel()
			errs = append(errs, err)
		}

	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func defaultServiceName(s IService) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath() + "." + t.Name()
}
