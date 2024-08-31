package services

import (
	"context"
	"fmt"

	"github.com/xarest/gobs"
)

type S1 struct {
	s2 *S2
	s3 *S3
}

var _ gobs.IServiceInit = (*S1)(nil)

func (s *S1) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("S1 init")
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S2), new(S3)},
	}, nil
}

var _ gobs.IServiceSetup = (*S1)(nil)

func (s *S1) Setup(ctx context.Context, deps gobs.Dependencies) error {
	fmt.Println("S1 setup")
	if err := deps.Assign(&s.s2, &s.s3); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	return nil
}

var _ gobs.IServiceStart = (*S1)(nil)

func (s *S1) Start(ctx context.Context) error {
	fmt.Println("S1 start")
	return nil
}

var _ gobs.IServiceStop = (*S1)(nil)

func (s *S1) Stop(ctx context.Context) error {
	fmt.Println("S1 stop")
	return nil
}
