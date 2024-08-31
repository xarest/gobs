package services

import (
	"context"
	"fmt"

	"github.com/xarest/gobs"
)

type S3 struct {
	s6 *S6
	s7 *S7
	s8 *S8
}

var _ gobs.IServiceInit = (*S1)(nil)

func (s *S3) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("S3 init")
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S6), new(S7), new(S8)},
	}, nil
}

var _ gobs.IServiceSetup = (*S1)(nil)

func (s *S3) Setup(ctx context.Context, deps gobs.Dependencies) error {
	fmt.Println("S3 setup")
	if err := deps.Assign(&s.s6, &s.s7, &s.s8); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	return nil
}

var _ gobs.IServiceStart = (*S1)(nil)

func (s *S3) Start(ctx context.Context) error {
	fmt.Println("S3 start")
	return nil
}

var _ gobs.IServiceStop = (*S1)(nil)

func (s *S3) Stop(ctx context.Context) error {
	fmt.Println("S3 stop")
	return nil
}
