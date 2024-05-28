package gobs

import "context"

type ServiceLifeCycle struct {
	OnSetup   *func(context.Context, []IService, []CustomService) error
	OnStart   *func(context.Context) error
	OnStop    *func(context.Context) error
	Deps      []IService
	ExtraDeps []CustomService
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

const (
	StatusInit  ServiceStatus = 0
	StatusSetup ServiceStatus = 1
	StatusStart ServiceStatus = 2
	StatusStop  ServiceStatus = 3
)
