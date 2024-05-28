package gobs

import (
	"context"
	"sync"
)

type Component struct {
	ServiceLifeCycle
	dependOn   []*Component
	followers  []*Component
	followChan []chan error
	service    IService
	name       string
	status     ServiceStatus
	mu         *sync.Mutex
	config     *Config
}

func NewComponent(s IService, name string, status ServiceStatus) *Component {
	return &Component{
		service: s,
		name:    name,
		status:  status,
		mu:      &sync.Mutex{},
	}
}

func (sb *Component) setup(ctx context.Context) (err error) {
	defer func() {
		sb.mu.Lock()
		sb.notify(err)
		sb.mu.Unlock()
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
	err = (*sb.OnSetup)(ctx, sb.Deps, sb.ExtraDeps)
	sb.mu.Lock()
	if err == nil {
		sb.status = StatusSetup
	}
	sb.mu.Unlock()
	return err
}

func (sb *Component) start(ctx context.Context) (err error) {
	defer func() {
		sb.mu.Lock()
		sb.notify(err)
		sb.mu.Unlock()
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

func (sb *Component) stop(ctx context.Context) (err error) {
	defer func() {
		sb.mu.Lock()
		sb.status = StatusStop
		sb.notify(err)
		sb.mu.Unlock()
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
		s.LogComponent("Service %s is waiting for %s", s.name, sb.name)
		sb.mu.Lock()
		if sb.status >= ss {
			sb.mu.Unlock()
			s.LogComponent("Service %s is done waiting for %s", s.name, sb.name)
			continue
		}
		ch := make(chan error)
		sb.followChan = append(sb.followChan, ch)
		sb.mu.Unlock()
		err = <-ch
		s.LogComponent("Service %s is done waiting for %s", s.name, sb.name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sb *Component) notify(err error) {
	sb.LogComponent("Service %s is notifying %d followers", sb.name, len(sb.followChan))
	for _, ch := range sb.followChan {
		ch <- err
		close(ch)
	}
	sb.followChan = nil
}
