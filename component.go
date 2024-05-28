package gobs

import (
	"context"
	"sync"
)

type BlockIdentifier struct {
	S IService
	N string
}

type Component struct {
	OnSetup      *func(context.Context, []BlockIdentifier) error
	OnStart      *func(context.Context) error
	OnStop       *func(context.Context) error
	Dependencies []BlockIdentifier

	dependOn   []*Component
	followers  []*Component
	followChan []chan error
	service    IService
	name       string
	status     ServiceStatus
	mu         *sync.Mutex
}

func NewComponent(s IService, name string, status ServiceStatus) *Component {
	return &Component{
		service: s,
		name:    name,
		status:  status,
		mu:      &sync.Mutex{},
	}
}

func (sb *Component) setup(ctx context.Context,
	onSetup *func(ctx context.Context, s IService, key string, err error),
) (err error) {
	defer func() {
		sb.mu.Lock()
		sb.notify(err)
		sb.mu.Unlock()
		(*onSetup)(ctx, sb.service, sb.name, err)
	}()
	if err = sb.wait(sb.dependOn, StatusSetup); err != nil {
		return err
	}
	sb.mu.Lock()
	if sb.status >= StatusSetup || sb.OnSetup == nil {
		sb.status = StatusSetup
		sb.mu.Unlock()
		return nil
	}
	sb.mu.Unlock()
	err = (*sb.OnSetup)(ctx, sb.Dependencies)
	sb.mu.Lock()
	if err == nil {
		sb.status = StatusSetup
	}
	sb.mu.Unlock()
	return err
}

func (sb *Component) start(ctx context.Context,
	onStart *func(ctx context.Context, s IService, key string, err error),
) (err error) {
	defer func() {
		sb.mu.Lock()
		sb.notify(err)
		sb.mu.Unlock()
		(*onStart)(ctx, sb.service, sb.name, err)
	}()
	if err = sb.wait(sb.dependOn, StatusStart); err != nil {
		return err
	}
	sb.mu.Lock()
	if sb.status == StatusInit || sb.status >= StatusStart || sb.OnStart == nil {
		sb.status = StatusStart
		sb.mu.Unlock()
		return nil
	}
	sb.mu.Unlock()
	err = (*sb.OnStart)(ctx)
	sb.mu.Lock()
	if err == nil {
		sb.status = StatusStart
	}
	sb.mu.Unlock()
	return err
}

func (sb *Component) stop(ctx context.Context,
	onStop *func(ctx context.Context, s IService, key string, err error),
) (err error) {
	defer func() {
		sb.mu.Lock()
		sb.status = StatusStop
		sb.notify(err)
		sb.mu.Unlock()
		if onStop != nil {
			(*onStop)(ctx, sb.service, sb.name, err)
		}
	}()
	if err = sb.wait(sb.followers, StatusStop); err != nil {
		return err
	}
	if sb.status == StatusInit || sb.status >= StatusStop || sb.OnStop == nil {
		return nil
	}
	return (*sb.OnStop)(ctx)
}

func (s *Component) wait(sbs []*Component, ss ServiceStatus) (err error) {
	for _, sb := range sbs {
		sb.mu.Lock()
		if sb.status >= ss {
			sb.mu.Unlock()
			continue
		}
		ch := make(chan error)
		sb.followChan = append(sb.followChan, ch)
		sb.mu.Unlock()

		err = <-ch
		if err != nil {
			return err
		}
	}
	return nil
}

func (sb *Component) notify(err error) {
	for _, ch := range sb.followChan {
		ch <- err
		close(ch)
	}
	sb.followChan = nil
}
