package gobs

import "context"

type Dependency struct {
	DependOn []IService
}

type IService interface {
	Init(context.Context, *Component) error
}

type ServiceStatus int

const (
	StatusInit  ServiceStatus = 0
	StatusSetup ServiceStatus = 1
	StatusStart ServiceStatus = 2
	StatusStop  ServiceStatus = 3
)
