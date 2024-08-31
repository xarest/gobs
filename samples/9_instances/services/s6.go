package services

import (
	"context"
	"fmt"
	"time"

	"github.com/xarest/gobs"
	"github.com/xarest/gobs/common"
)

type S6 struct {
	s10 *S10
	s11 *S11
}

var _ gobs.IServiceInit = (*S6)(nil)

func (s *S6) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("S6 init")
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S10), new(S11)},
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
			common.StatusStop:  true,
		},
	}, nil
}

var _ gobs.IServiceSetup = (*S6)(nil)

func (s *S6) Setup(ctx context.Context, deps gobs.Dependencies) error {
	time.Sleep(30 * time.Millisecond)
	fmt.Println("S6 setup")
	if err := deps.Assign(&s.s10, &s.s11); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	return nil
}

var _ gobs.IServiceStart = (*S6)(nil)

func (s *S6) Start(ctx context.Context) error {
	fmt.Println("S6 start")
	return nil
}

var _ gobs.IServiceStop = (*S6)(nil)

func (s *S6) Stop(ctx context.Context) error {
	time.Sleep(30 * time.Millisecond)
	fmt.Println("S6 stop")
	return nil
}
