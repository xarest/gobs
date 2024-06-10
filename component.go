package gobs

import (
	"context"
	"sync"

	"github.com/traphamxuan/gobs/logger"
)

type Component struct {
	ServiceLifeCycle
	*logger.Logger
	following []*Component
	followers []*Component
	service   IService
	name      string
	status    ServiceStatus
	mu        *sync.Mutex
	// config    *Config
}

func NewComponent(s IService, name string, status ServiceStatus, log *logger.Logger) *Component {
	c := &Component{
		Logger:  log,
		service: s,
		name:    name,
		status:  status,
		mu:      &sync.Mutex{},
	}
	c.AddTag("Component/" + name)
	return c
}

func (sb *Component) ShouldRun(ss ServiceStatus) (err error) {
	untag := sb.AddTag("ShouldRun")
	defer untag()
	if sb.status >= ss {
		sb.Log("Skip service %s because it already ran", sb.name)
		return ErrorServiceRan
	}
	if ss >= StatusStop && sb.status < StatusSetup {
		sb.Log("Skip service %s because it has not been setup", sb.name)
		sb.mu.Lock()
		sb.status = StatusStop
		sb.mu.Unlock()
		return ErrorServiceRan
	}
	dependOn := sb.DependOn(ss)
	for _, dep := range dependOn {
		if dep.status < ss {
			sb.Log("Skip service %s because it depends on %s", sb.name, dep.name)
			return ErrorServiceNotReady
		}
	}
	return nil
}

func (sb *Component) DependOn(ss ServiceStatus) []*Component {
	if ss >= StatusStop {
		return sb.followers
	}
	return sb.following
}

func (sb *Component) Followers(ss ServiceStatus) []*Component {
	if ss >= StatusStop {
		return sb.following
	}
	return sb.followers
}

func (sb *Component) Run(ctx context.Context, ss ServiceStatus) (err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	if sb.status >= ss {
		sb.Log("Skip service %s because it already ran", sb.name)
		return ErrorServiceRan
	}
	switch ss {
	case StatusSetup:
		if sb.OnSetup != nil {
			err = (*sb.OnSetup)(ctx, sb.Deps, sb.ExtraDeps)
		} else if sb.OnSetupAsync != nil {
			err = (*sb.OnSetupAsync)(ctx, sb.Deps, sb.ExtraDeps)
		}
	case StatusStart:
		if sb.OnStart != nil {
			err = (*sb.OnStart)(ctx)
		}
	case StatusStop:
		if sb.OnStop != nil {
			err = (*sb.OnStop)(ctx)
		} else if sb.OnStopAsync != nil {
			err = (*sb.OnStopAsync)(ctx)
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

func (sb *Component) IsRunAsync(ss ServiceStatus) bool {
	switch ss {
	case StatusSetup:
		return sb.OnSetupAsync != nil
	case StatusStart:
		return sb.OnStart != nil
	case StatusStop:
		return sb.OnStopAsync != nil
	}
	return false
}
