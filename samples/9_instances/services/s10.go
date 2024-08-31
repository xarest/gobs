package services

import (
	"context"
	"fmt"
	"time"

	"github.com/xarest/gobs"
	"github.com/xarest/gobs/common"
)

type S10 struct{}

var _ gobs.IServiceInit = (*S10)(nil)

func (s *S10) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("S10 init")
	return &gobs.ServiceLifeCycle{
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
			common.StatusStop:  true,
		},
	}, nil
}

var _ gobs.IServiceSetup = (*S10)(nil)

func (s *S10) Setup(ctx context.Context, deps gobs.Dependencies) error {
	time.Sleep(60 * time.Millisecond)
	fmt.Println("S10 setup")

	return nil
}

var _ gobs.IServiceStart = (*S10)(nil)

func (s *S10) Start(ctx context.Context) error {
	fmt.Println("S10 start")
	return nil
}

var _ gobs.IServiceStop = (*S10)(nil)

func (s *S10) Stop(ctx context.Context) error {
	time.Sleep(60 * time.Millisecond)
	fmt.Println("S10 stop")
	return nil
}
