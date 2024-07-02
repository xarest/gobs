package gobs

import (
	"context"
	"errors"

	"github.com/traphamxuan/gobs/common"
	"github.com/traphamxuan/gobs/logger"
	"github.com/traphamxuan/gobs/scheduler"
	"github.com/traphamxuan/gobs/types"
	"github.com/traphamxuan/gobs/utils"
)

type Bootstrap struct {
	*logger.Logger
	numOfConcurrencies int
	schedulers         map[common.ServiceStatus]*scheduler.Scheduler
	services           []*Service
	keys               map[string]*Service
}

func NewBootstrap(configs ...Config) *Bootstrap {
	cfg := DefaultConfig
	if len(configs) > 0 {
		cfg = configs[0]
	}

	bs := &Bootstrap{
		Logger: logger.NewLog(cfg.Logger),
		// config: cfg,
		schedulers:         make(map[common.ServiceStatus]*scheduler.Scheduler, common.StatusStop+1),
		numOfConcurrencies: cfg.NumOfConcurrencies,
		keys:               make(map[string]*Service),
	}
	bs.SetDetail(cfg.EnableLogDetail)
	bs.SetTag("Bootstrap")
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
			return cp.instance
		}
		return nil
	}
	key = utils.DefaultServiceName(service)
	bs.Log("Get service with key %s", key)
	if cp, ok := bs.keys[key]; ok {
		return cp.instance
	}
	return nil
}

func (bs *Bootstrap) AddDefault(s IService, args ...string) error {
	if len(args) > 0 {
		return bs.Add(s, common.StatusNone, args[0])
	}
	return bs.Add(s, common.StatusNone, "")
}

func (bs *Bootstrap) AddOrPanic(s IService, args ...string) {
	if err := bs.AddDefault(s, args...); err != nil {
		panic(err)
	}
}

func (bs *Bootstrap) Add(s IService, status common.ServiceStatus, key string) error {
	untag := bs.AddTag("Add")
	defer untag()
	if key == "" {
		key = utils.DefaultServiceName(s)
		if key == "" {
			return errors.New("service name is empty")
		}
	}

	if bs.keys[key] != nil {
		return nil
	}
	sBlock := NewService(s, key, status, bs.Logger.Clone())
	bs.keys[key] = sBlock
	bs.services = append(bs.services, sBlock)
	bs.LogS("Service %s is added with status %s", utils.CompactName(key), status.String())
	return nil
}

func (bs *Bootstrap) Init(ctx context.Context) error {
	untag := bs.AddTag("Init")
	defer untag()
	totalLength := len(bs.services)
	var tasks []types.ITask
	for i := 0; i < totalLength; i++ {
		sb := bs.services[i]
		unTag := bs.AddTag(utils.CompactName(sb.name))
		if err := sb.instance.Init(ctx, sb); err != nil {
			bs.LogS("Failed to init %s", sb.name, err.Error())
			return err
		}
		bs.LogS("Initialized successfully")

		if err := bs.setupNetworkConnection(sb); err != nil {
			bs.LogS("Failed to set dependencies, %s", sb.name, err.Error())
			return err
		}
		tasks = append(tasks, sb)
		totalLength = len(bs.services)
		unTag()
	}
	return bs.execute(ctx, common.StatusInit, tasks)
}

func (bs *Bootstrap) Setup(ctx context.Context) error {
	sched, ok := bs.schedulers[common.StatusInit]
	if !ok {
		return errors.New("Init is not executed")
	}
	sched.Interrupt()
	tasks, err := sched.Release()
	if err != nil {
		bs.Logger.Log("Previous state %s hash error %s", common.StatusInit.String(), err.Error())
	}
	untag := bs.AddTag("Setup")
	defer untag()
	return bs.execute(ctx, common.StatusSetup, tasks)
}

func (bs *Bootstrap) Start(ctx context.Context) error {
	sched, ok := bs.schedulers[common.StatusSetup]
	if !ok {
		return errors.New("Setup is not executed")
	}
	sched.Interrupt()
	tasks, err := sched.Release()
	if err != nil {
		bs.Logger.Log("Previous state %s hash error %s", common.StatusSetup.String(), err.Error())
	}
	untag := bs.AddTag("Start")
	defer untag()
	return bs.execute(ctx, common.StatusStart, tasks)
}

func (bs *Bootstrap) Stop(ctx context.Context) error {
	sched, ok := bs.schedulers[common.StatusStart]
	if ok && sched != nil {
		sched.Interrupt()
	}
	sched, ok = bs.schedulers[common.StatusSetup]
	if !ok || sched == nil {
		bs.LogS("Setup is not executed. Skip stop")
		return nil
	}
	sched.Interrupt()
	tasks, err := sched.Release()
	if err != nil {
		bs.Logger.Log("Previous state %s hash error %s", common.StatusSetup.String(), err.Error())
	}
	untag := bs.AddTag("Stop")
	defer untag()
	sched = scheduler.NewScheduler(ctx, bs.Logger.Clone(), tasks, common.StatusStop, bs.numOfConcurrencies)
	for _, service := range bs.services {
		if service.status < common.StatusSetup {
			sched.SetIgnore(service)
		}
	}
	bs.schedulers[common.StatusStop] = sched
	return sched.Run(ctx)
}

func (bs *Bootstrap) Break(ctx context.Context) {
	for k := range bs.schedulers {
		bs.schedulers[k].Interrupt()
	}
}

func (bs *Bootstrap) execute(ctx context.Context, ss common.ServiceStatus, tasks []types.ITask) (err error) {
	untag := bs.AddTag("execute-" + ss.String())
	defer untag()
	bs.LogS("Execute %s with %d tasks", ss.String(), len(tasks))
	sched := scheduler.NewScheduler(ctx, bs.Logger.Clone(), tasks, ss, bs.numOfConcurrencies)
	bs.schedulers[ss] = sched
	return sched.Run(ctx)
}

func (bs *Bootstrap) setupNetworkConnection(sb *Service) error {
	var services []IService
	for _, service := range sb.Deps {
		key := utils.DefaultServiceName(service)

		dService, ok := bs.keys[key]
		if !ok {
			bs.Add(service, common.StatusNone, key)
			dService = bs.keys[key]
		}
		sb.following = append(sb.following, dService)
		dService.followers = append(dService.followers, sb)
		services = append(services, dService.instance)
	}
	sb.Deps = services

	var extraServices []CustomService
	for _, cService := range sb.ExtraDeps {
		key := cService.Name
		if key == "" {
			key = utils.DefaultServiceName(cService.Service)
		}
		dService, ok := bs.keys[key]
		if !ok {
			if cService.Instance != nil {
				iService, ok := cService.Instance.(IService)
				if ok {
					bs.Add(iService, common.StatusNone, key)
				} else {
					bs.Add(cService.Service, common.StatusNone, key)
				}
			} else {
				bs.Add(cService.Service, common.StatusNone, key)
			}
			dService = bs.keys[key]
		}
		sb.following = append(sb.following, dService)
		dService.followers = append(dService.followers, sb)
		extraServices = append(extraServices, CustomService{
			Service:  dService.instance,
			Name:     key,
			Instance: cService.Instance,
		})
	}
	sb.ExtraDeps = extraServices
	return nil
}
