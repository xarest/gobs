package services

import (
	"context"
	"fmt"

	"github.com/xarest/gobs"
)

type S12 struct{}

var _ gobs.IServiceInit = (*S12)(nil)

func (s *S12) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("S12 init")
	return nil, nil
}

var _ gobs.IServiceSetup = (*S12)(nil)

func (s *S12) Setup(ctx context.Context, deps gobs.Dependencies) error {
	fmt.Println("S12 setup")

	return nil
}

var _ gobs.IServiceStart = (*S12)(nil)

func (s *S12) Start(ctx context.Context) error {
	fmt.Println("S12 start")
	return nil
}

var _ gobs.IServiceStop = (*S12)(nil)

func (s *S12) Stop(ctx context.Context) error {
	fmt.Println("S12 stop")
	return nil
}
