package gobs

import (
	"context"
	"errors"

	"github.com/traphamxuan/gobs/logger"
)

type Bootstrap struct {
	*logger.Logger
	// config    Config
	scheduler *Scheduler
	services  []*Component
	keys      map[string]*Component
}

func NewBootstrap(configs ...Config) *Bootstrap {
	cfg := DefaultConfig
	if len(configs) > 0 {
		cfg = configs[0]
	}

	bs := &Bootstrap{
		Logger: logger.NewLog(cfg.Logger),
		// config: cfg,
		keys: make(map[string]*Component),
	}
	bs.SetDetail(cfg.EnableLogDetail)
	bs.SetTag("Bootstrap")
	bs.scheduler = NewScheduler(cfg.NumOfConcurrencies, bs.Logger.Clone())
	return bs
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
	untag := bs.AddTag("Add")
	defer untag()
	if key == "" {
		key = defaultServiceName(s)
		if key == "" {
			return errors.New("service name is empty")
		}
	}

	if bs.keys[key] != nil {
		return nil
	}
	sBlock := NewComponent(s, key, status, bs.Logger.Clone())
	bs.keys[key] = sBlock
	bs.services = append(bs.services, sBlock)
	bs.LogS("Service %s is added", key)
	return nil
}

func (bs *Bootstrap) Init(ctx context.Context) error {
	untag := bs.AddTag("Init")
	defer untag()
	totalLength := len(bs.services)
	for i := 0; i < totalLength; i++ {
		sb := bs.services[i]
		if err := sb.service.Init(ctx, sb); err != nil {
			return err
		}

		if err := bs.setupNetworkConnection(sb); err != nil {
			return err
		}
		totalLength = len(bs.services)
		bs.Log("New length after init %d, %p", totalLength, sb.OnSetupAsync)
	}
	bs.scheduler.Load(bs.services)
	return nil
}

func (bs *Bootstrap) Setup(ctx context.Context) error {
	return bs.scheduler.Run(ctx, StatusSetup)
}

func (bs *Bootstrap) Start(ctx context.Context) error {
	return bs.scheduler.Run(ctx, StatusStart)
}

func (bs *Bootstrap) Stop(ctx context.Context) error {
	return bs.scheduler.Run(ctx, StatusStop)
}

func (bs *Bootstrap) setupNetworkConnection(sb *Component) error {
	var services []IService
	for _, service := range sb.Deps {
		key := defaultServiceName(service)

		dComponent, ok := bs.keys[key]
		if !ok {
			bs.Add(service, StatusInit, key)
			dComponent = bs.keys[key]
		}
		sb.following = append(sb.following, dComponent)
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
			if cService.Instance != nil {
				iService, ok := cService.Instance.(IService)
				if ok {
					bs.Add(iService, StatusInit, key)
				} else {
					bs.Add(cService.Service, StatusInit, key)
				}
			} else {
				bs.Add(cService.Service, StatusInit, key)
			}
			dComponent = bs.keys[key]
		}
		sb.following = append(sb.following, dComponent)
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
