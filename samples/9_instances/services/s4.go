package services

import (
	"context"
	"fmt"

	"github.com/xarest/gobs"
)

type S4 struct {
	s9  *S9
	s10 *S10
}

var _ gobs.IServiceInit = (*S4)(nil)

func (s *S4) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("S4 init")
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S9), new(S10)},
	}, nil
}

var _ gobs.IServiceSetup = (*S4)(nil)

func (s *S4) Setup(ctx context.Context, deps gobs.Dependencies) error {
	fmt.Println("S4 setup")
	if err := deps.Assign(&s.s9, &s.s10); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	return nil
}

var _ gobs.IServiceStart = (*S4)(nil)

func (s *S4) Start(ctx context.Context) error {
	fmt.Println("S4 start")
	return nil
}

var _ gobs.IServiceStop = (*S4)(nil)

func (s *S4) Stop(ctx context.Context) error {
	fmt.Println("S4 stop")
	return nil
}
