package gobs

import (
	"context"
	"sync"

	"github.com/xarest/gobs/common"
	"github.com/xarest/gobs/logger"
	"github.com/xarest/gobs/types"
	"github.com/xarest/gobs/utils"
)

type IService any

type IServiceInit interface {
	// Entry point to connect service intances with the others. This method will be called at the beginning of the bootstrap process
	// to build up the dependencies between services. This method will setup `s *Service` lifecycle.
	//
	// Example:
	//
	// func (d *D) Init(ctx context.Context, s *gobs.Service) error {
	// 	s.Deps = []gobs.IService{&B{}, &C{}} // Define dependencies here
	// 	s.OnSetup = func(ctx context.Context, deps []gobs.IService, extraDeps []gobs.CustomService) error {
	// 		// After B & C finish setting up, this callback will be called
	// 		d.B = deps[0].(*B)
	// 		d.C = deps[1].(*C)
	// 		// Other custom setup/configration go here
	// 		return nil
	// 	}
	// 	s.AsyncMode[common.StatusSetup] = true // This line will make OnSetup method be called in concurrent context without blocking others.
	// 	return nil
	// }
	Init(ctx context.Context) (*ServiceLifeCycle, error)
}

// Without configuration of OnSetup at Init(...), this method will be called when main process invokes bootstrap.Setup(...).
type IServiceSetup interface {
	Setup(ctx context.Context, deps ...IService) error
}

// Start method will be called when the main context invokes the bootstrap.Start(...) method and OnStart method is not configured at Init(...).
type IServiceStart interface {
	Start(ctx context.Context) error
}

type IServiceStartServer interface {
	StartServer(ctx context.Context, onReady func(err error)) error
}

// Stop method will be called when the main context invokes the bootstrap.Stop(...) method and OnStop method is not configured at Init(...).
type IServiceStop interface {
	Stop(ctx context.Context) error
}

type ServiceLifeCycle struct {
	OnInterrupt func(errno int)

	// deprecated: Define func () Setup(ctx context.Context, deps Dependencies) error instead of OnSetup.
	// OnSetup is a callback function that will be called when the main context invokes the bootstrap.Setup(...) method.
	// This method is used to assign the dependencies instances which has setup successfully from gobs to the service instance.
	// The `deps` parameter is a list of dependencies that the service instance depends on. The `extraDeps` parameter is a list of
	// custom dependencies in case service don't share dependencies with the others.
	AfterInit func(ctx context.Context, deps ...IService) error

	// Deps is a list of dependencies that the service instance depends on. This list is just a reference to the type struct.
	// Gobs will automatically look up the existing instances or create new instances based on the type struct.
	// Then set dependencies to the service instance after Init(...) method returns nil.
	Deps Dependencies

	// ExtraDeps is a list of custom dependencies that the service instance depends on.
	// It provides more information about the instance that the service instance depends on.
	//
	// Example:
	//
	// func (d *D) Init(ctx context.Context, s *gobs.Service) error {
	// 	s.ExtraDeps = []gobs.CustomService{
	// 		{&B{}, "", instanceB}, // D depends on instanceB, type struct B{} with key is default
	// 		{&C{}, "C1", nil}, // D depends on an instance C which has key is C1 in gobs
	// 	}
	// }
	ExtraDeps []CustomService

	// AsyncMode is a map of service status and boolean value. If the value is true, the service instance will be run in parallel goroutine context.
	// Otherwise, the service instance will be run in sequential context.
	AsyncMode map[common.ServiceStatus]bool
}

type CustomService struct {
	Service  IService
	Name     string
	Instance IService
}

type Service struct {
	ServiceLifeCycle
	*logger.Logger
	following []types.ITask
	followers []types.ITask
	instance  IService
	name      string
	status    common.ServiceStatus
	mutex     map[common.ServiceStatus]*sync.Mutex
}

var _ types.ITask = (*Service)(nil)

func (sb *Service) Name() string {
	return sb.name
}

func NewService(s any, name string, status common.ServiceStatus, log *logger.Logger) *Service {
	c := &Service{
		ServiceLifeCycle: ServiceLifeCycle{
			AsyncMode: make(map[common.ServiceStatus]bool, common.StatusStop+1),
		},
		Logger:   log,
		instance: s,
		name:     name,
		status:   status,
		mutex: map[common.ServiceStatus]*sync.Mutex{
			common.StatusUninitialized: {},
			common.StatusInit:          {},
			common.StatusSetup:         {},
			common.StatusStart:         {},
			common.StatusStop:          {},
		},
	}
	c.AddTag("Service/" + name)
	return c
}

func (sb *Service) DependOn(ss common.ServiceStatus) []types.ITask {
	if ss >= common.StatusStop {
		return sb.followers
	}
	return sb.following
}

func (sb *Service) Followers(ss common.ServiceStatus) []types.ITask {
	if ss >= common.StatusStop {
		return sb.following
	}
	return sb.followers
}

func (sb *Service) Run(ctx context.Context, ss common.ServiceStatus) (err error) {
	mutex, ok := sb.mutex[ss]
	if ok && mutex != nil {
		mutex.Lock()
		defer mutex.Unlock()
	}
	logKey := utils.CompactName(sb.name)
	switch ss {
	case common.StatusInit:
		if sb.AfterInit != nil {
			err = sb.AfterInit(ctx, sb.Deps...)
		}
	case common.StatusSetup:
		if s, ok := sb.instance.(IServiceSetup); ok {
			err = s.Setup(ctx, sb.Deps...)
		} else {
			sb.Log("Service %s does not implement IServiceSetup", logKey)
		}
	case common.StatusStart:
		if s, ok := sb.instance.(IServiceStart); ok {
			err = s.Start(ctx)
		} else if s, ok := sb.instance.(IServiceStartServer); ok {
			var chErr = make(chan error)
			go func(ctx context.Context) {
				defer close(chErr)
				s.StartServer(ctx, func(e error) {
					chErr <- e
				})
			}(ctx)
			err = <-chErr
		} else {
			sb.Log("Service %s does not implement IServiceStart", logKey)
		}
	case common.StatusStop:
		if s, ok := sb.instance.(IServiceStop); ok {
			err = s.Stop(ctx)
		} else {
			sb.Log("Service %s does not implement IServiceStop", logKey)
		}
	default:
		err = nil
	}

	if err != nil {
		return err
	}
	sb.status = ss
	return nil
}

func (sb *Service) IsRunAsync(ss common.ServiceStatus) bool {
	return sb.AsyncMode[ss]
}

func (sb *Service) UpdateDependencies(dep *Service) {
	sb.following = append(sb.following, dep)
	dep.followers = append(dep.followers, sb)
	logKey := utils.CompactName(sb.name)
	logServiceKey := utils.CompactName(dep.name)
	sb.Log("Service %s depends on service %s (%d) and service %s has service %s as a follower (%d)",
		logServiceKey, logKey, len(sb.following), logKey, logServiceKey, len(dep.followers),
	)
}
