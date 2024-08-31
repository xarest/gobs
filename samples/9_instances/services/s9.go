package services

import (
	"context"
	"fmt"
	"time"

	"github.com/xarest/gobs"
	"github.com/xarest/gobs/common"
)

type S9 struct{}

var _ gobs.IServiceInit = (*S9)(nil)

func (s *S9) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("S9 init")
	return &gobs.ServiceLifeCycle{
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
			common.StatusStop:  true,
		},
	}, nil
}

var _ gobs.IServiceSetup = (*S9)(nil)

func (s *S9) Setup(ctx context.Context, deps gobs.Dependencies) error {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("S9 setup")

	return nil
}

var _ gobs.IServiceStart = (*S9)(nil)

func (s *S9) Start(ctx context.Context) error {
	fmt.Println("S9 start")
	return nil
}

var _ gobs.IServiceStop = (*S9)(nil)

func (s *S9) Stop(ctx context.Context) error {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("S9 stop")
	return nil
}
