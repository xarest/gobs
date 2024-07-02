package gobs_test

import (
	"context"
	"time"

	"github.com/traphamxuan/gobs"
	"github.com/traphamxuan/gobs/common"
)

func commonSetup(id int, err error, delayms time.Duration) func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
	return func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		time.Sleep(delayms * time.Millisecond)
		setupOrder = append(setupOrder, id)
		return err
	}
}
func commonStop(id int, err error, delayms time.Duration) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		time.Sleep(delayms * time.Millisecond)
		stopOrder = append(stopOrder, id)
		return err
	}
}

var setupOrder = []int{}
var stopOrder = []int{}

type S1 struct {
	err error
	S2  *S2
	S3  *S3
}

func (s *S1) Init(ctx context.Context, co *gobs.Service) error {
	co.Deps = []gobs.IService{new(S2), new(S3)}
	co.OnStop = commonStop(1, nil, 0)
	co.OnSetup = func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S2 = deps[0].(*S2)
		s.S3 = deps[1].(*S3)
		setupOrder = append(setupOrder, 1)
		return s.err
	}
	return nil
}

var _ gobs.IService = (*S1)(nil)

type S2 struct {
	err error
	S4  *S4
	S5  *S5
}

func (s *S2) Init(ctx context.Context, co *gobs.Service) error {
	co.Deps = []gobs.IService{new(S4), new(S5)}
	co.OnStop = commonStop(2, nil, 0)
	co.OnSetup = func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S4 = deps[0].(*S4)
		s.S5 = deps[1].(*S5)
		setupOrder = append(setupOrder, 2)
		return s.err
	}
	return nil
}

var _ gobs.IService = (*S2)(nil)

type S3 struct {
	err error
	S6  *S6
	S7  *S7
	S8  *S8
}

func (s *S3) Init(ctx context.Context, co *gobs.Service) error {
	co.Deps = []gobs.IService{new(S6), new(S7), new(S8)}
	co.OnStop = commonStop(3, nil, 0)
	co.OnSetup = func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S6 = deps[0].(*S6)
		s.S7 = deps[1].(*S7)
		s.S8 = deps[2].(*S8)
		setupOrder = append(setupOrder, 3)
		return s.err
	}
	return nil
}

var _ gobs.IService = (*S3)(nil)

type S4 struct {
	err error
	S9  *S9
	S10 *S10
}

func (s *S4) Init(ctx context.Context, co *gobs.Service) error {
	co.Deps = []gobs.IService{new(S9), new(S10)}
	co.OnStop = commonStop(4, nil, 0)
	co.OnSetup = func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S9 = deps[0].(*S9)
		s.S10 = deps[1].(*S10)
		setupOrder = append(setupOrder, 4)
		return s.err
	}
	return nil
}

var _ gobs.IService = (*S4)(nil)

type S5 struct {
	err error
	S9  *S9
	S10 *S10
	S11 *S11
}

func (s *S5) Init(ctx context.Context, co *gobs.Service) error {
	co.Deps = []gobs.IService{new(S9), new(S10), new(S11)}
	co.OnStop = commonStop(5, nil, 0)
	co.OnSetup = func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S9 = deps[0].(*S9)
		s.S10 = deps[1].(*S10)
		s.S11 = deps[2].(*S11)
		setupOrder = append(setupOrder, 5)
		return s.err
	}
	return nil
}

var _ gobs.IService = (*S5)(nil)

type S6 struct {
	err error
	S10 *S10
	S11 *S11
}

func (s *S6) Init(ctx context.Context, co *gobs.Service) error {
	co.Deps = []gobs.IService{new(S10), new(S11)}
	co.OnStop = commonStop(6, nil, 0)
	co.OnSetup = func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S10 = deps[0].(*S10)
		s.S11 = deps[1].(*S11)
		setupOrder = append(setupOrder, 6)
		return s.err
	}
	return nil
}

var _ gobs.IService = (*S6)(nil)

type S7 struct {
	err error
	S12 *S12
}

func (s *S7) Init(ctx context.Context, co *gobs.Service) error {
	co.Deps = []gobs.IService{new(S12)}
	co.OnStop = commonStop(7, nil, 0)
	co.OnSetup = func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S12 = deps[0].(*S12)
		setupOrder = append(setupOrder, 7)
		return s.err
	}
	return nil
}

var _ gobs.IService = (*S7)(nil)

type S8 struct {
	err error
	S13 *S13
}

func (s *S8) Init(ctx context.Context, co *gobs.Service) error {
	co.Deps = []gobs.IService{new(S13)}
	co.OnStop = commonStop(8, nil, 0)
	co.OnSetup = func(ctx context.Context, deps []gobs.IService, _ []gobs.CustomService) error {
		s.S13 = deps[0].(*S13)
		setupOrder = append(setupOrder, 8)
		return s.err
	}
	return nil
}

var _ gobs.IService = (*S8)(nil)

type S9 struct{ err error }

func (s *S9) Init(ctx context.Context, co *gobs.Service) error {
	co.OnStop = commonStop(9, nil, 0)
	co.OnSetup = commonSetup(9, s.err, 100)
	co.AsyncMode[common.StatusSetup] = true
	return nil
}

var _ gobs.IService = (*S9)(nil)

type S10 struct{ err error }

func (s *S10) Init(ctx context.Context, co *gobs.Service) error {
	co.OnStop = commonStop(10, nil, 0)
	co.OnSetup = commonSetup(10, s.err, 90)
	co.AsyncMode[common.StatusSetup] = true
	return nil
}

var _ gobs.IService = (*S10)(nil)

type S11 struct{ err error }

func (s *S11) Init(ctx context.Context, co *gobs.Service) error {
	co.OnStop = commonStop(11, nil, 0)
	co.OnSetup = commonSetup(11, s.err, 1)
	return nil
}

var _ gobs.IService = (*S11)(nil)

type S12 struct{ err error }

func (s *S12) Init(ctx context.Context, co *gobs.Service) error {
	co.OnStop = commonStop(12, nil, 0)
	co.OnSetup = commonSetup(12, s.err, 2)
	return nil
}

var _ gobs.IService = (*S12)(nil)

type S13 struct{ err error }

func (s *S13) Init(ctx context.Context, co *gobs.Service) error {
	co.OnStop = commonStop(13, nil, 0)
	co.OnSetup = commonSetup(13, s.err, 80)
	co.AsyncMode[common.StatusSetup] = true
	return nil
}

var _ gobs.IService = (*S13)(nil)
