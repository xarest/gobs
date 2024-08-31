package services

import (
	"context"
	"fmt"

	"github.com/xarest/gobs"
)

type S11 struct{}

var _ gobs.IServiceInit = (*S11)(nil)

func (s *S11) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("S11 init")
	return nil, nil
}

var _ gobs.IServiceSetup = (*S11)(nil)

func (s *S11) Setup(ctx context.Context, deps gobs.Dependencies) error {
	fmt.Println("S11 setup")

	return nil
}

var _ gobs.IServiceStart = (*S11)(nil)

func (s *S11) Start(ctx context.Context) error {
	fmt.Println("S11 start")
	return nil
}

var _ gobs.IServiceStop = (*S11)(nil)

func (s *S11) Stop(ctx context.Context) error {
	fmt.Println("S11 stop")
	return nil
}
