package gobs

import (
	"context"
	"errors"
)

type Bootstrap struct {
	config   Config
	services []*Component
	keys     map[string]*Component
}

func NewBootstrap(configs ...Config) *Bootstrap {
	cfg := DefaultConfig
	if len(configs) > 0 {
		cfg = configs[0]
	}
	return &Bootstrap{
		config: cfg,
		keys:   make(map[string]*Component),
	}
}

func (bs *Bootstrap) Deinit(ctx context.Context) {
	bs.keys = nil
	bs.services = nil
}

func (bs *Bootstrap) GetService(service IService, key string) IService {
	if service == nil {
		if key == "" {
			return nil
		}
		if cp, ok := bs.keys[key]; ok {
			return cp.service
		}
		return nil
	}
	key = defaultServiceName(service)
	if cp, ok := bs.keys[key]; ok {
		return cp.service
	}
	return nil
}

func (bs *Bootstrap) AddDefault(s IService, args ...string) error {
	if len(args) > 0 {
		return bs.Add(s, StatusInit, args[0])
	}
	return bs.Add(s, StatusInit, "")
}

func (bs *Bootstrap) AddOrPanic(s IService, args ...string) {
	if err := bs.AddDefault(s, args...); err != nil {
		panic(err)
	}
}

func (bs *Bootstrap) Add(s IService, status ServiceStatus, key string) error {
	if key == "" {
		key = defaultServiceName(s)
	}
	if key == "" {
		return errors.New("service name is empty")
	}
	if bs.keys[key] != nil {
		return nil
	}
	sBlock := NewComponent(s, key, status)
	bs.keys[key] = sBlock
	bs.services = append(bs.services, sBlock)
	sBlock.config = &bs.config
	bs.LogModule(bs.config.EnableLogAdd, "Service %s is added", key)
	return nil
}

func (bs *Bootstrap) Setup(ctx context.Context) error {
	for i := 0; i < len(bs.services); i++ {
		sb := bs.services[i]
		if err := sb.service.Init(ctx, sb); err != nil {
			return err
		}
		if err := bs.dependenciesToFollowers(sb); err != nil {
			return err
		}
	}
	return concurrenceProcesses(ctx, bs.services,
		func(ctx context.Context, sb *Component) error {
			err := sb.setup(ctx)
			bs.LogModule(bs.config.EnableLogSetup, "Service %s setup successfully", sb.name)
			return err
		},
	)
}

func (bs *Bootstrap) dependenciesToFollowers(sb *Component) error {
	var services []IService
	for _, service := range sb.Deps {
		key := defaultServiceName(service)

		dComponent, ok := bs.keys[key]
		if !ok {
			bs.Add(service, StatusInit, key)
			dComponent = bs.keys[key]
		}
		sb.dependOn = append(sb.dependOn, dComponent)
		dComponent.followers = append(dComponent.followers, sb)
		services = append(services, dComponent.service)
	}
	sb.Deps = services

	var extraServices []CustomService
	for _, cService := range sb.ExtraDeps {
		key := cService.Name
		if key == "" {
			key = defaultServiceName(cService.Service)
		}
		dComponent, ok := bs.keys[key]
		if !ok {
			bs.Add(cService.Service, StatusInit, key)
			dComponent = bs.keys[key]
		}
		sb.dependOn = append(sb.dependOn, dComponent)
		dComponent.followers = append(dComponent.followers, sb)
		extraServices = append(extraServices, CustomService{
			Service:  dComponent.service,
			Name:     key,
			Instance: cService.Instance,
		})
	}
	sb.ExtraDeps = extraServices
	return nil
}

func (bs *Bootstrap) Start(ctx context.Context) error {
	return concurrenceProcesses(ctx, bs.services,
		func(ctx context.Context, sb *Component) error {
			err := sb.start(ctx)
			bs.LogModule(bs.config.EnableLogStart, "Service %s started", sb.name)
			return err
		},
	)
}

func (bs *Bootstrap) Stop(ctx context.Context) error {
	return concurrenceProcesses(ctx, bs.services,
		func(ctx context.Context, sb *Component) error {
			sb.stop(ctx)
			bs.LogModule(bs.config.EnableLogStop, "Service %s is stopped", sb.name)
			return nil
		},
	)
}
