package gobs_test

import (
	"context"

	"github.com/traphamxuan/gobs"
)

type A struct{}

func (a *A) Init(ctx context.Context, co *gobs.Service) error {
	co.OnSetup = func(ctx context.Context, deps []gobs.IService, extraDeps []gobs.CustomService) error {
		return nil
	}
	return nil
}

var _ gobs.IService = (*B)(nil)

type B struct {
	A *A
}

func (b *B) Init(ctx context.Context, co *gobs.Service) error {
	co.Deps = []gobs.IService{&A{}}
	co.OnSetup = func(ctx context.Context, deps []gobs.IService, extraDeps []gobs.CustomService) error {
		b.A = deps[0].(*A)
		return nil
	}
	return nil
}

var _ gobs.IService = (*C)(nil)

type C struct {
	A *A
	B *B
}

func (c *C) Init(ctx context.Context, co *gobs.Service) error {
	co.Deps = []gobs.IService{&A{}, &B{}}
	co.OnSetup = func(ctx context.Context, deps []gobs.IService, extraDeps []gobs.CustomService) error {
		c.A = deps[0].(*A)
		c.B = deps[1].(*B)
		return nil
	}
	return nil
}

var _ gobs.IService = (*D)(nil)

type D struct {
	B *B
	C *C
}

func (d *D) Init(ctx context.Context, co *gobs.Service) error {
	co.Deps = []gobs.IService{&B{}, &C{}}
	co.OnSetup = func(ctx context.Context, deps []gobs.IService, extraDeps []gobs.CustomService) error {
		d.B = deps[0].(*B)
		d.C = deps[1].(*C)
		return nil
	}
	return nil
}
