package services

import (
	"context"
	"fmt"
	"time"

	"github.com/xarest/gobs"
	"github.com/xarest/gobs/common"
)

type S13 struct{}

var _ gobs.IServiceInit = (*S13)(nil)

func (s *S13) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("S13 init")
	return &gobs.ServiceLifeCycle{
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
			common.StatusStop:  true,
		},
	}, nil
}

var _ gobs.IServiceSetup = (*S13)(nil)

func (s *S13) Setup(ctx context.Context, deps gobs.Dependencies) error {
	time.Sleep(30 * time.Millisecond)
	fmt.Println("S13 setup")

	return nil
}

var _ gobs.IServiceStart = (*S13)(nil)

func (s *S13) Start(ctx context.Context) error {
	fmt.Println("S13 start")
	return nil
}

var _ gobs.IServiceStop = (*S13)(nil)

func (s *S13) Stop(ctx context.Context) error {
	fmt.Println("S13 stop")
	return nil
}
