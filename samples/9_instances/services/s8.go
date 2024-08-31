package services

import (
	"context"
	"fmt"

	"github.com/xarest/gobs"
)

type S8 struct {
	s13 *S13
}

var _ gobs.IServiceInit = (*S8)(nil)

func (s *S8) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("S8 init")
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S13)},
	}, nil
}

var _ gobs.IServiceSetup = (*S8)(nil)

func (s *S8) Setup(ctx context.Context, deps gobs.Dependencies) error {
	fmt.Println("S8 setup")
	if err := deps.Assign(&s.s13); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	return nil
}

var _ gobs.IServiceStart = (*S8)(nil)

func (s *S8) Start(ctx context.Context) error {
	fmt.Println("S8 start")
	return nil
}

var _ gobs.IServiceStop = (*S8)(nil)

func (s *S8) Stop(ctx context.Context) error {
	fmt.Println("S8 stop")
	return nil
}
