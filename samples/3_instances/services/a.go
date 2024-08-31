package services

import (
	"context"
	"fmt"

	"github.com/xarest/gobs"
)

type A struct{}

var _ gobs.IServiceInit = (*A)(nil)

func (s A) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("A Init")
	return nil, nil
}

var _ gobs.IServiceSetup = (*A)(nil)

func (s A) Setup(ctx context.Context, deps gobs.Dependencies) error {
	fmt.Println("A Setup")
	return nil
}

var _ gobs.IServiceStart = (*A)(nil)

func (s A) Start(ctx context.Context) error {
	fmt.Println("A Start")
	return nil
}

var _ gobs.IServiceStop = (*A)(nil)

func (s A) Stop(ctx context.Context) error {
	fmt.Println("A Stop")
	return nil
}
