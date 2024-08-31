package services

import (
	"context"
	"fmt"

	"github.com/xarest/gobs"
)

type S2 struct {
	s4 *S4
	s5 *S5
}

var _ gobs.IServiceInit = (*S1)(nil)

func (s *S2) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("S2 init")
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S4), new(S5)},
	}, nil
}

var _ gobs.IServiceSetup = (*S1)(nil)

func (s *S2) Setup(ctx context.Context, deps gobs.Dependencies) error {
	fmt.Println("S2 setup")
	if err := deps.Assign(&s.s4, &s.s5); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	return nil
}

var _ gobs.IServiceStart = (*S1)(nil)

func (s *S2) Start(ctx context.Context) error {
	fmt.Println("S2 start")
	return nil
}

var _ gobs.IServiceStop = (*S1)(nil)

func (s *S2) Stop(ctx context.Context) error {
	fmt.Println("S2 stop")
	return nil
}
