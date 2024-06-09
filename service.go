package gobs

import "context"

type ServiceLifeCycle struct {
	OnSetup      *func(context.Context, []IService, []CustomService) error
	OnSetupAsync *func(context.Context, []IService, []CustomService) error
	OnStart      *func(context.Context) error
	OnStop       *func(context.Context) error
	OnStopAsync  *func(context.Context) error
	Deps         []IService
	ExtraDeps    []CustomService
}

type IService interface {
	Init(context.Context, *Component) error
}

type CustomService struct {
	Service  IService
	Name     string
	Instance interface{}
}

type ServiceStatus int

func (ss ServiceStatus) String() string {
	switch ss {
	case StatusInit:
		return "Init"
	case StatusSetup:
		return "Setup"
	case StatusStart:
		return "Start"
	case StatusStop:
		return "Stop"
	default:
		return "Unknown"
	}
}

const (
	StatusInit  ServiceStatus = 0
	StatusSetup ServiceStatus = 1
	StatusStart ServiceStatus = 2
	StatusStop  ServiceStatus = 3
)
