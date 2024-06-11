package gobs_test

import (
	"context"

	"github.com/traphamxuan/gobs"
)

func commonSetup(id int) func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
	return func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		setupOrder = append(setupOrder, id)
		return nil
	}
}

var setupOrder = []int{}

type S1 struct {
	S2 *S2
	S3 *S3
}

func (s *S1) Init(ctx context.Context, co *gobs.Component) error {
	co.Deps = []gobs.IService{new(S2), new(S3)}
	onSetup := func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S2 = deps[0].(*S2)
		s.S3 = deps[1].(*S3)
		setupOrder = append(setupOrder, 1)
		return nil
	}
	co.OnSetup = &onSetup
	return nil
}

var _ gobs.IService = (*S1)(nil)

type S2 struct {
	S4 *S4
	S5 *S5
}

func (s *S2) Init(ctx context.Context, co *gobs.Component) error {
	co.Deps = []gobs.IService{new(S4), new(S5)}
	onSetup := func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S4 = deps[0].(*S4)
		s.S5 = deps[1].(*S5)
		setupOrder = append(setupOrder, 2)
		return nil
	}
	co.OnSetup = &onSetup
	return nil
}

var _ gobs.IService = (*S2)(nil)

type S3 struct {
	S6 *S6
	S7 *S7
	S8 *S8
}

func (s *S3) Init(ctx context.Context, co *gobs.Component) error {
	co.Deps = []gobs.IService{new(S6), new(S7), new(S8)}
	onSetup := func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S6 = deps[0].(*S6)
		s.S7 = deps[1].(*S7)
		s.S8 = deps[2].(*S8)
		setupOrder = append(setupOrder, 3)
		return nil
	}
	co.OnSetup = &onSetup
	return nil
}

var _ gobs.IService = (*S3)(nil)

type S4 struct {
	S9  *S9
	S10 *S10
}

func (s *S4) Init(ctx context.Context, co *gobs.Component) error {
	co.Deps = []gobs.IService{new(S9), new(S10)}
	onSetup := func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S9 = deps[0].(*S9)
		s.S10 = deps[1].(*S10)
		setupOrder = append(setupOrder, 4)
		return nil
	}
	co.OnSetup = &onSetup
	return nil
}

var _ gobs.IService = (*S4)(nil)

type S5 struct {
	S9  *S9
	S10 *S10
	S11 *S11
}

func (s *S5) Init(ctx context.Context, co *gobs.Component) error {
	co.Deps = []gobs.IService{new(S9), new(S10), new(S11)}
	onSetup := func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S9 = deps[0].(*S9)
		s.S10 = deps[1].(*S10)
		s.S11 = deps[2].(*S11)
		setupOrder = append(setupOrder, 5)
		return nil
	}
	co.OnSetup = &onSetup
	return nil
}

var _ gobs.IService = (*S5)(nil)

type S6 struct {
	S10 *S10
	S11 *S11
}

func (s *S6) Init(ctx context.Context, co *gobs.Component) error {
	co.Deps = []gobs.IService{new(S10), new(S11)}
	onSetup := func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S10 = deps[0].(*S10)
		s.S11 = deps[1].(*S11)
		setupOrder = append(setupOrder, 6)
		return nil
	}
	co.OnSetup = &onSetup
	return nil
}

var _ gobs.IService = (*S6)(nil)

type S7 struct {
	S12 *S12
}

func (s *S7) Init(ctx context.Context, co *gobs.Component) error {
	co.Deps = []gobs.IService{new(S12)}
	onSetup := func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S12 = deps[0].(*S12)
		setupOrder = append(setupOrder, 7)
		return nil
	}
	co.OnSetup = &onSetup
	return nil
}

var _ gobs.IService = (*S7)(nil)

type S8 struct {
	S13 *S13
}

func (s *S8) Init(ctx context.Context, co *gobs.Component) error {
	co.Deps = []gobs.IService{new(S13)}
	onSetup := func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S13 = deps[0].(*S13)
		setupOrder = append(setupOrder, 8)
		return nil
	}
	co.OnSetup = &onSetup
	return nil
}

var _ gobs.IService = (*S8)(nil)

type S9 struct{}

func (s *S9) Init(ctx context.Context, co *gobs.Component) error {
	setupFnc := commonSetup(9)
	co.OnSetupAsync = &setupFnc
	return nil
}

var _ gobs.IService = (*S9)(nil)

type S10 struct{}

func (s *S10) Init(ctx context.Context, co *gobs.Component) error {
	setupFnc := commonSetup(10)
	co.OnSetupAsync = &setupFnc
	return nil
}

var _ gobs.IService = (*S10)(nil)

type S11 struct{}

func (s *S11) Init(ctx context.Context, co *gobs.Component) error {
	setupFnc := commonSetup(11)
	co.OnSetup = &setupFnc
	return nil
}

var _ gobs.IService = (*S11)(nil)

type S12 struct{}

func (s *S12) Init(ctx context.Context, co *gobs.Component) error {
	setupFnc := commonSetup(12)
	co.OnSetup = &setupFnc
	return nil
}

var _ gobs.IService = (*S12)(nil)

type S13 struct{}

func (s *S13) Init(ctx context.Context, co *gobs.Component) error {
	setupFnc := commonSetup(13)
	co.OnSetupAsync = &setupFnc
	return nil
}

var _ gobs.IService = (*S13)(nil)
