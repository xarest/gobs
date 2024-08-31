package services

import (
	"context"
	"fmt"

	"github.com/xarest/gobs"
)

type S5 struct {
	s9  *S9
	s10 *S10
	s11 *S11
}

var _ gobs.IServiceInit = (*S5)(nil)

func (s *S5) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("S5 init")
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S9), new(S10), new(S11)},
	}, nil
}

var _ gobs.IServiceSetup = (*S5)(nil)

func (s *S5) Setup(ctx context.Context, deps gobs.Dependencies) error {
	fmt.Println("S5 setup")
	if err := deps.Assign(&s.s9, &s.s10, &s.s11); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	return nil
}

var _ gobs.IServiceStart = (*S5)(nil)

func (s *S5) Start(ctx context.Context) error {
	fmt.Println("S5 start")
	return nil
}

var _ gobs.IServiceStop = (*S5)(nil)

func (s *S5) Stop(ctx context.Context) error {
	fmt.Println("S5 stop")
	return nil
}
