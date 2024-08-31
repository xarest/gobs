package services

import (
	"context"
	"fmt"

	"github.com/xarest/gobs"
)

type C struct{}

var _ gobs.IServiceInit = (*C)(nil)

func (s C) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("C Init")
	return nil, nil
}

var _ gobs.IServiceSetup = (*C)(nil)

func (s C) Setup(ctx context.Context, deps gobs.Dependencies) error {
	fmt.Println("C Setup")
	return nil
}

var _ gobs.IServiceStart = (*C)(nil)

func (s C) Start(ctx context.Context) error {
	fmt.Println("C Start")
	return nil
}

var _ gobs.IServiceStop = (*C)(nil)

func (s C) Stop(ctx context.Context) error {
	fmt.Println("C Stop")
	return nil
}
