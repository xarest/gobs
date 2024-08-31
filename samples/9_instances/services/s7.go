package services

import (
	"context"
	"fmt"
	"time"

	"github.com/xarest/gobs"
	"github.com/xarest/gobs/common"
)

type S7 struct {
	s12 *S12
}

var _ gobs.IServiceInit = (*S7)(nil)

func (s *S7) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("S7 init")
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S12)},
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
			common.StatusStop:  true,
		},
	}, nil
}

var _ gobs.IServiceSetup = (*S7)(nil)

func (s *S7) Setup(ctx context.Context, deps gobs.Dependencies) error {
	time.Sleep(50 * time.Millisecond)
	fmt.Println("S7 setup")
	if err := deps.Assign(&s.s12); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	return nil
}

var _ gobs.IServiceStart = (*S7)(nil)

func (s *S7) Start(ctx context.Context) error {
	fmt.Println("S7 start")
	return nil
}

var _ gobs.IServiceStop = (*S7)(nil)

func (s *S7) Stop(ctx context.Context) error {
	time.Sleep(50 * time.Millisecond)
	fmt.Println("S7 stop")
	return nil
}
