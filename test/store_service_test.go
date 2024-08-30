package gobs_test

import (
	"context"

	"github.com/xarest/gobs"
)

type A struct{}

func (a *A) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return nil, nil
}

var _ gobs.IService = (*B)(nil)

type B struct {
	A *A
}

func (b *B) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: []gobs.IService{&A{}},
		OnSetup: func(ctx context.Context, deps gobs.Dependencies) error {
			return deps.Assign(&b.A)
		},
	}, nil
}

var _ gobs.IService = (*C)(nil)

type C struct {
	A *A
	B *B
}

func (c *C) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: []gobs.IService{&A{}, &B{}},
		OnSetup: func(ctx context.Context, deps gobs.Dependencies) error {
			return deps.Assign(&c.A, &c.B)
		},
	}, nil
}

var _ gobs.IService = (*D)(nil)

type D struct {
	B *B
	C *C
}

func (d *D) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: []gobs.IService{&B{}, &C{}},
		OnSetup: func(ctx context.Context, deps gobs.Dependencies) error {
			deps.Assign(&d.B, &d.C)
			return nil
		},
	}, nil
}
