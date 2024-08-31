package services

import (
	"context"
	"fmt"

	"github.com/xarest/gobs"
)

type B struct{}

var _ gobs.IServiceInit = (*B)(nil)

func (s B) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	fmt.Println("B Init")
	return nil, nil
}

var _ gobs.IServiceSetup = (*B)(nil)

func (s B) Setup(ctx context.Context, deps gobs.Dependencies) error {
	fmt.Println("B Setup")
	return nil
}

var _ gobs.IServiceStart = (*B)(nil)

func (s B) Start(ctx context.Context) error {
	fmt.Println("B Start")
	return nil
}

var _ gobs.IServiceStop = (*B)(nil)

func (s B) Stop(ctx context.Context) error {
	fmt.Println("B Stop")
	return nil
}
