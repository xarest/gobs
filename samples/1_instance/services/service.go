package services

import (
	"context"
	"fmt"

	"github.com/xarest/gobs"
)

type Service struct{}

var _ gobs.IServiceInit = (*Service)(nil)

func (s Service) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("Service Init")
	return nil, nil
}

var _ gobs.IServiceSetup = (*Service)(nil)

func (s Service) Setup(ctx context.Context, deps gobs.Dependencies) error {
	fmt.Println("Service Setup")
	return nil
}

var _ gobs.IServiceStart = (*Service)(nil)

func (s Service) Start(ctx context.Context) error {
	fmt.Println("Service Start")
	return nil
}

var _ gobs.IServiceStop = (*Service)(nil)

func (s Service) Stop(ctx context.Context) error {
	fmt.Println("Service Stop")
	return nil
}
