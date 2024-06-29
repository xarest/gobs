package gobs

import (
	"context"
	"sync"

	"github.com/traphamxuan/gobs/common"
	"github.com/traphamxuan/gobs/logger"
	"github.com/traphamxuan/gobs/types"
	"github.com/traphamxuan/gobs/utils"
)

type IService interface {
	Init(context.Context, *Service) error
}

type ServiceLifeCycle struct {
	OnSetup   func(context.Context, []IService, []CustomService) error
	OnStart   func(context.Context) error
	OnStop    func(context.Context) error
	Deps      []IService
	ExtraDeps []CustomService
	AsyncMode map[common.ServiceStatus]bool
}

type CustomService struct {
	Service  IService
	Name     string
	Instance interface{}
}

type Service struct {
	ServiceLifeCycle
	*logger.Logger
	following []types.ITask
	followers []types.ITask
	instance  IService
	name      string
	status    common.ServiceStatus
	mu        *sync.Mutex
	// config    *Config
}

var _ types.ITask = (*Service)(nil)

func (sb *Service) Name() string {
	return sb.name
}

func NewService(s IService, name string, status common.ServiceStatus, log *logger.Logger) *Service {
	c := &Service{
		ServiceLifeCycle: ServiceLifeCycle{
			AsyncMode: make(map[common.ServiceStatus]bool, common.StatusStop+1),
			OnSetup:   func(_ context.Context, _ []IService, _ []CustomService) error { return nil },
			OnStart:   utils.EmptyFunc,
			OnStop:    utils.EmptyFunc,
		},
		Logger:   log,
		instance: s,
		name:     name,
		status:   status,
		mu:       &sync.Mutex{},
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
	sb.mu.Lock()
	defer sb.mu.Unlock()
	if sb.status >= ss {
		sb.Log("Skip service %s at %s because it is %s", sb.name, ss.String(), sb.status.String())
		return utils.ErrorServiceRan
	}
	switch ss {
	case common.StatusSetup:
		err = sb.OnSetup(ctx, sb.Deps, sb.ExtraDeps)
	case common.StatusStart:
		err = sb.OnStart(ctx)
	case common.StatusStop:
		err = sb.OnStop(ctx)
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
