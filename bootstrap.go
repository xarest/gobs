package gobs

import (
	"context"
	"fmt"
	"reflect"
)

type Bootstrap struct {
	services []*Component
	keys     map[string]*Component
}

func NewBootstrap() *Bootstrap {
	return &Bootstrap{
		keys: make(map[string]*Component),
	}
}

func (sm *Bootstrap) Deinit(ctx context.Context) {
	sm.keys = nil
	sm.services = nil
}

// Please add a unique name for your instance so that other can use it to find
// Without a unique name `sm.AddDefault(&Router{})`, sm will use the type name `Router`
// If you have multiple instance with the same name, they will conflicts
// Ex: middleware.Authentication vs service.Authentication
func (sm *Bootstrap) AddDefault(s IService, args ...string) {
	sm.Add(s, StatusInit, args...)
}

func (sm *Bootstrap) Add(s IService, status ServiceStatus, args ...string) {
	key := reflect.TypeOf(s).Elem().Name()
	if len(args) > 0 && args[0] != "" {
		key = args[0]
	}
	if key == "" {
		panic("Service key is empty")
	}
	if sm.keys[key] != nil {
		return
	}
	sBlock := NewComponent(s, key, status)
	sm.keys[key] = sBlock
	sm.services = append(sm.services, sBlock)
}

func (sm *Bootstrap) Setup(
	ctx context.Context,
	onSetup *func(ctx context.Context, s IService, key string, err error),
) error {
	for _, sb := range sm.services {
		if err := sb.service.Init(ctx, sb); err != nil {
			return err
		}
		if err := sm.dependenciesToFollowers(sb); err != nil {
			return err
		}
	}
	return ConcurrenceProcesses(ctx, sm.services,
		func(ctx context.Context, sb *Component) error {
			return sb.setup(ctx, onSetup)
		},
	)
}

func (sm *Bootstrap) dependenciesToFollowers(sb *Component) error {
	var dependencies []BlockIdentifier
	for _, dep := range sb.Dependencies {
		key := reflect.TypeOf(dep.S).Elem().Name()
		if dep.N != "" {
			key = dep.N
		}
		depend, ok := sm.keys[key]
		if !ok {
			return fmt.Errorf("service %s depends on %s, but %s is not found, %w", sb.name, key, key, ErrorServiceNotFound)
		}
		sb.dependOn = append(sb.dependOn, depend)
		depend.followers = append(depend.followers, sb)
		dependencies = append(dependencies, BlockIdentifier{
			S: depend.service,
			N: depend.name,
		})
	}
	sb.Dependencies = dependencies
	return nil
}

func (sm *Bootstrap) Start(
	ctx context.Context,
	onStart *func(ctx context.Context, s IService, key string, err error),
) error {
	return ConcurrenceProcesses(ctx, sm.services,
		func(ctx context.Context, sb *Component) error {
			return sb.start(ctx, onStart)
		},
	)
}

func (sm *Bootstrap) Stop(
	ctx context.Context,
	onStop *func(ctx context.Context, s IService, key string, err error),
) error {
	return ConcurrenceProcesses(ctx, sm.services,
		func(ctx context.Context, sb *Component) error {
			// fmt.Println("Stopping service", sb.name)
			sb.stop(ctx, onStop)
			// fmt.Println("Stopped service", sb.name)
			return nil
		},
	)
}
