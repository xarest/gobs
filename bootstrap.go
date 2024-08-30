package gobs

import (
	"context"
	"errors"
	"time"

	"github.com/xarest/gobs/common"
	"github.com/xarest/gobs/logger"
	"github.com/xarest/gobs/scheduler"
	"github.com/xarest/gobs/types"
	"github.com/xarest/gobs/utils"
)

type Bootstrap struct {
	*logger.Logger
	numOfConcurrencies int
	schedulers         map[common.ServiceStatus]*scheduler.Scheduler
	services           []*Service
	keys               map[string]*Service
}

// NewBootstrap creates a new Bootstrap instance using the provided configurations.
// If no configuration options are passed, it applies the DefaultConfig
//
// Example:
//
//	import (
//		"github.com/xarest/gobs"
//		"github.com/xarest/gobs/logger"
//		"fmt"
//	)
//
//	var log logger.LogFnc = func(s string, i ...interface{}) {
//		fmt.Printf(s+"\n", i...)
//	}
//
//	bs := gobs.NewBootstrap(gobs.Config {
//		NumOfConcurrencies: -1, // -1 means unlimited
//		EnableLogDetail:    false,
//		Logger:             &log,
//	})
func NewBootstrap(configs ...Config) *Bootstrap {
	cfg := DefaultConfig
	if len(configs) > 0 {
		cfg = configs[0]
	}

	bs := &Bootstrap{
		Logger:             logger.NewLog(cfg.Logger),
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

// GetService returns the service instance with the provided original instance or key.
// If the service is not found, it returns nil. Key or instance can be nil or empty but not both.
// If both are provided, the key will be preferred.
//
// Example:
//
//	a := bs.GetService(nil, "abc") // Return instance of service which has key "abc"
//	a := bs.GetService(&A{}, "") // Return instance of service type A which has default key which is path of struct A
func GetService[T IService](bs *Bootstrap, service T, key string) (*T, bool) {
	if key == "" {
		key = utils.DefaultServiceName(service)
	}
	if cp, ok := bs.keys[key]; ok {
		res, ok := cp.instance.(*T)
		return res, ok
	}
	return &service, false
}

// AddDefault is wrapper for Add method with default key and status of the service instance.
// If key is not provided, it will use the default key which is the path of the service instance.
// If the service is already added, no duplicate instances for a key are allowed.
// If the service instance is added with no key and no path (nearly impossible), it will return an error
//
// Example:
//
//	bs.AddDefault(new(A)) // Add service instance A with default key which is path of struct A
//	bs.AddDefault(new(A), "abc") // Add another service instance A with key "abc".
func (bs *Bootstrap) AddDefault(s IService, args ...string) error {
	if len(args) > 0 {
		return bs.Add(s, common.StatusUninitialized, args[0])
	}
	return bs.Add(s, common.StatusUninitialized, "")
}

// Same with AddDefault but panic if error
func (bs *Bootstrap) AddOrPanic(s IService, args ...string) {
	if err := bs.AddDefault(s, args...); err != nil {
		panic(err)
	}
}

func (bs *Bootstrap) AddMany(services ...IService) error {
	for _, s := range services {
		if err := bs.AddDefault(s); err != nil {
			return err
		}
	}
	return nil
}

// Add method is used to add a service instance to the bootstrap.
// The service instance must implement IService interface.
// The status of the service instance can be any value from common.ServiceStatus.
// It is helpful but not recommend to setup service instance before adding to the bootstrap.
//
// Example:
//
//	func main() {
//		log = &logger.Logger{}
//
//		if err := api.log.Setup(ctx); err != nil {
//			panic(err)
//		}
//
//		var l gl.LogFnc = api.log.Debugf
//
//		bs := gobs.NewBootstrap(gobs.Config{
//			NumOfConcurrencies: -1,
//			Logger:             &l,
//		})
//
//		bs.Add(log, common.StatusSetup, "") // log instance will skip setup process
//	}
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

// Init method is used to initialize all services in the bootstrap.
// It must be called before Setup to build dependencies between services.
// All Init(...) method of services implemented IService interface will be called.
// Those methods in services will be called in sequence.
func (bs *Bootstrap) Init(ctx context.Context) error {
	untag := bs.AddTag("Init")
	defer untag()
	totalLength := len(bs.services)
	var tasks []types.ITask
	bs.Log("INIT WITH %d SERVICES", totalLength)
	for i := 0; i < totalLength; i++ {
		sb := bs.services[i]
		taskKey := utils.CompactName(sb.name)
		unTag := bs.AddTag(taskKey)
		if inst, ok := sb.instance.(ServiceInit); ok {
			sCfg, err := inst.Init(ctx)
			if err != nil {
				bs.LogS("Failed to init %s", taskKey, err.Error())
				return err
			}
			if sCfg != nil {
				if err := bs.setupNetworkConnection(sb, *sCfg); err != nil {
					bs.LogS("Failed to set dependencies, %s", sb.name, err.Error())
					return err
				}
			}
		}

		tasks = append(tasks, sb)
		totalLength = len(bs.services)
		unTag()
	}
	return bs.execute(ctx, common.StatusInit, tasks, 0)
}

// Setup method is used to setup all services in the bootstrap.
// It must be called before Start(...) method. Results of setup process (internnally) will be used in Start(...) method.
// Make sure that the Init(...) method is fisnished before calling this method. Otherwise, it will interrupt the Init(...) process
// and return an error as Init(...) method is not finished.
func (bs *Bootstrap) Setup(ctx context.Context) error {
	sched, ok := bs.schedulers[common.StatusInit]
	if !ok {
		return errors.New("Init is not executed")
	}
	sched.Interrupt()
	tasks, err := sched.Release()
	if err != nil {
		bs.Logger.LogS("Previous state %s has error %s", common.StatusInit.String(), err.Error())
	}
	untag := bs.AddTag("Setup")
	defer untag()
	return bs.execute(ctx, common.StatusSetup, tasks, bs.numOfConcurrencies)
}

// Start method is used to start all services in the bootstrap.
// It must be called after Setup(...) method. Results of setup process (internnally) will be used in Start(...) method.
// Make sure that the Setup(...) method is fisnished before calling this method. Otherwise, it will interrupt the Setup(...) process
// and return an error as Setup(...) method is not finished.
// If service instances use Start() method as container for holding pending states that OnStart function cannot return,
// it won't be marked as started.
// If you set other services depended on pending services, make sure that the pending service has it own goroutine
// to handle the pending states and OnStart function must return. It's not recommended to use Start() method for this purpose.
func (bs *Bootstrap) Start(ctx context.Context) error {
	sched, ok := bs.schedulers[common.StatusSetup]
	if !ok {
		return errors.New("Setup is not executed")
	}
	sched.Interrupt()
	tasks, err := sched.Release()
	if err != nil {
		bs.Logger.LogS("Previous state %s has error %s", common.StatusSetup.String(), err.Error())
	}
	untag := bs.AddTag("Start")
	defer untag()
	return bs.execute(ctx, common.StatusStart, tasks, bs.numOfConcurrencies)
}

// Stop method is used to stop all services which have been setup successfully.
// This method will try to interrupt all pending states of services in Start(...) method and wait for them to finish. Before invoking OnStop method.
// Stop method is the must-have method to call before the application is terminated. Its flows are inverted of Setup(...) method.
// If service B depends on services A, service A will be stopped after service B is stopped.
func (bs *Bootstrap) Stop(ctx context.Context) error {
	sched, ok := bs.schedulers[common.StatusStart]
	if ok && sched != nil {
		sched.Interrupt()
	}
	sched, ok = bs.schedulers[common.StatusSetup]
	if !ok || sched == nil {
		bs.LogS("Setup is not executed. Skip stopping process")
		return nil
	}
	sched.Interrupt()
	tasks, err := sched.Release()
	if err != nil {
		bs.Logger.LogS("Previous state %s has error %s", common.StatusSetup.String(), err.Error())
	}
	untag := bs.AddTag("Stop")
	defer untag()
	bs.LogS("EXECUTE %s WITH %d SERVICES", common.StatusStop.String(), len(tasks))
	sched = scheduler.NewScheduler(ctx, bs.Logger.Clone(), tasks, common.StatusStop, bs.numOfConcurrencies)
	for _, service := range bs.services {
		if service.status < common.StatusSetup {
			sched.SetIgnore(service)
		}
	}
	bs.schedulers[common.StatusStop] = sched
	return sched.Run(ctx)
}

// Interrupt method is used to notify all services to stopo their processes.
// If a service are waiting for other services to finish, it will be interrupted and stop waiting.
// If a service are running, it will continute to run until it finishes.
//
// For example:
// router are waiting for success of database connection, if interrupt is called, router will stop waiting for database connection and return without setting up.
// if router are running, it will continue to serve requests until OnStop(...) was called to safely shutdown router.
func (bs *Bootstrap) Interrupt(ctx context.Context) {
	for k := range bs.schedulers {
		bs.schedulers[k].Interrupt()
	}
}

func (bs *Bootstrap) StartBootstrap(ctx context.Context, quits ...chan struct{}) {
	appCtx, cancelAll := context.WithCancel(ctx)
	defer cancelAll()
	ctxDone, cancel := context.WithCancel(appCtx)
	defer cancel()
	go func(ctx context.Context, cancel context.CancelFunc) {
		defer cancel()
		if err := bs.Init(ctx); err != nil {
			bs.LogS("Failed to init services: %s", err.Error())
			return
		}
		if err := bs.Setup(ctx); err != nil {
			bs.LogS("Failed to setup services: %s", err.Error())
			return
		}
		bs.Start(ctx)
	}(appCtx, cancel)

	utils.WaitOnEvents(ctxDone, func(ctx context.Context, event struct{}) error {
		return common.ErrorEndOfProcessing
	}, nil, quits...)
	bs.Interrupt(ctxDone)

	quitCtx, done := context.WithTimeout(appCtx, 10*time.Second)
	defer done()
	go func() {
		defer done()
		bs.Stop(quitCtx)
		bs.Deinit(quitCtx)
	}()
	// catching ctx.Done(). timeout of 5 seconds.
	<-quitCtx.Done()
}

func (bs *Bootstrap) execute(ctx context.Context, ss common.ServiceStatus, tasks []types.ITask, numOfConcurrencies int) (err error) {
	untag := bs.AddTag("execute-" + ss.String())
	defer untag()
	bs.LogS("EXECUTE %s WITH %d SERVICES", ss.String(), len(tasks))
	sched := scheduler.NewScheduler(ctx, bs.Logger.Clone(), tasks, ss, numOfConcurrencies)
	bs.schedulers[ss] = sched
	return sched.Run(ctx)
}

func (bs *Bootstrap) setupNetworkConnection(sb *Service, sCfg ServiceLifeCycle) error {
	logServiceKey := utils.CompactName(sb.name)
	bs.Log("Service %s has %d dependencies", logServiceKey, len(sCfg.Deps))

	for _, service := range sCfg.Deps {
		key := utils.DefaultServiceName(service)

		dService, ok := bs.keys[key]
		if !ok {
			if err := bs.Add(service, common.StatusUninitialized, key); err != nil {
				return err
			}
			dService = bs.keys[key]
		}
		sb.UpdateDependencies(dService)
	}

	bs.Log("Service %s has %d extra dependencies", logServiceKey, len(sCfg.ExtraDeps))
	for _, cService := range sCfg.ExtraDeps {
		key := cService.Name
		if key == "" {
			key = utils.DefaultServiceName(cService.Service)
		}
		dService, ok := bs.keys[key]
		if !ok || dService == nil {
			if cService.Instance != nil {
				iService, ok := cService.Instance.(IService)
				if ok {
					if err := bs.Add(iService, common.StatusUninitialized, key); err != nil {
						return err
					}
				} else {
					if err := bs.Add(cService.Service, common.StatusUninitialized, key); err != nil {
						return err
					}
				}
			} else {
				if err := bs.Add(cService.Service, common.StatusUninitialized, key); err != nil {
					return err
				}
			}
			dService = bs.keys[key]
		}
		sb.UpdateDependencies(dService)
	}
	sCfg.Deps = nil
	for _, dep := range sb.following {
		if d, ok := dep.(*Service); ok {
			sCfg.Deps = append(sCfg.Deps, d.instance)
		}
	}
	sb.ServiceLifeCycle = sCfg
	return nil
}
