package gobs_test

import (
	"context"
	"fmt"
	"time"

	"github.com/traphamxuan/gobs"
	"github.com/traphamxuan/gobs/common"
)

func commonSetup(id int, err error, delayms time.Duration) func(ctx context.Context, deps gobs.Dependencies) error {
	return func(ctx context.Context, deps gobs.Dependencies) error {
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

func (s *S1) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps:   []gobs.IService{new(S2), new(S3)},
		OnStop: commonStop(1, nil, 0),
		OnSetup: func(ctx context.Context, deps gobs.Dependencies) error {
			if err := deps.Assign(&s.S2, &s.S3); err != nil {
				fmt.Println("Failed to assign dependencies", err)
				return err
			}
			setupOrder = append(setupOrder, 1)
			return s.err
		},
	}, nil
}

var _ gobs.IService = (*S1)(nil)

type S2 struct {
	err error
	S4  *S4
	S5  *S5
}

func (s *S2) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps:   []gobs.IService{new(S4), new(S5)},
		OnStop: commonStop(2, nil, 0),
		OnSetup: func(ctx context.Context, deps gobs.Dependencies) error {
			if err := deps.Assign(&s.S4, &s.S5); err != nil {
				fmt.Println("Failed to assign dependencies", err)
				return err
			}
			setupOrder = append(setupOrder, 2)
			return s.err
		},
	}, nil
}

var _ gobs.IService = (*S2)(nil)

type S3 struct {
	err error
	S6  *S6
	S7  *S7
	S8  *S8
}

func (s *S3) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps:   []gobs.IService{new(S6), new(S7), new(S8)},
		OnStop: commonStop(3, nil, 0),
		OnSetup: func(ctx context.Context, deps gobs.Dependencies) error {
			if err := deps.Assign(&s.S6, &s.S7, &s.S8); err != nil {
				fmt.Println("Failed to assign dependencies", err)
				return err
			}
			setupOrder = append(setupOrder, 3)
			return s.err
		},
	}, nil
}

var _ gobs.IService = (*S3)(nil)

type S4 struct {
	err error
	S9  *S9
	S10 *S10
}

func (s *S4) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps:   []gobs.IService{new(S9), new(S10)},
		OnStop: commonStop(4, nil, 0),
		OnSetup: func(ctx context.Context, deps gobs.Dependencies) error {
			if err := deps.Assign(&s.S9, &s.S10); err != nil {
				fmt.Println("Failed to assign dependencies", err)
				return err
			}
			setupOrder = append(setupOrder, 4)
			return s.err
		},
	}, nil
}

var _ gobs.IService = (*S4)(nil)

type S5 struct {
	err error
	S9  *S9
	S10 *S10
	S11 *S11
}

func (s *S5) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps:   []gobs.IService{new(S9), new(S10), new(S11)},
		OnStop: commonStop(5, nil, 0),
		OnSetup: func(ctx context.Context, deps gobs.Dependencies) error {
			if err := deps.Assign(&s.S9, &s.S10, &s.S11); err != nil {
				fmt.Println("Failed to assign dependencies", err)
				return err
			}
			setupOrder = append(setupOrder, 5)
			return s.err
		},
	}, nil
}

var _ gobs.IService = (*S5)(nil)

type S6 struct {
	err error
	S10 *S10
	S11 *S11
}

func (s *S6) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps:   []gobs.IService{new(S10), new(S11)},
		OnStop: commonStop(6, nil, 0),
		OnSetup: func(ctx context.Context, deps gobs.Dependencies) error {
			if err := deps.Assign(&s.S10, &s.S11); err != nil {
				fmt.Println("Failed to assign dependencies", err)
				return err
			}
			setupOrder = append(setupOrder, 6)
			return s.err
		},
	}, nil
}

var _ gobs.IService = (*S6)(nil)

type S7 struct {
	err error
	S12 *S12
}

func (s *S7) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps:   []gobs.IService{new(S12)},
		OnStop: commonStop(7, nil, 0),
		OnSetup: func(ctx context.Context, deps gobs.Dependencies) error {
			if err := deps.Assign(&s.S12); err != nil {
				fmt.Println("Failed to assign dependencies", err)
				return err
			}
			setupOrder = append(setupOrder, 7)
			return s.err
		},
	}, nil
}

var _ gobs.IService = (*S7)(nil)

type S8 struct {
	err error
	S13 *S13
}

func (s *S8) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps:   []gobs.IService{new(S13)},
		OnStop: commonStop(8, nil, 0),
		OnSetup: func(ctx context.Context, deps gobs.Dependencies) error {
			if err := deps.Assign(&s.S13); err != nil {
				fmt.Println("Failed to assign dependencies", err)
				return err
			}
			setupOrder = append(setupOrder, 8)
			return s.err
		},
	}, nil
}

var _ gobs.IService = (*S8)(nil)

type S9 struct{ err error }

func (s *S9) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		OnStop:  commonStop(9, nil, 0),
		OnSetup: commonSetup(9, s.err, 100),
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
		},
	}, nil
}

var _ gobs.IService = (*S9)(nil)

type S10 struct{ err error }

func (s *S10) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
		},
		OnSetup: commonSetup(10, s.err, 90),
		OnStop:  commonStop(10, nil, 0),
	}, nil
}

var _ gobs.IService = (*S10)(nil)

type S11 struct{ err error }

func (s *S11) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
		},
		OnSetup: commonSetup(11, s.err, 1),
		OnStop:  commonStop(11, nil, 0),
	}, nil
}

var _ gobs.IService = (*S11)(nil)

type S12 struct{ err error }

func (s *S12) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
		},
		OnSetup: commonSetup(12, s.err, 2),
		OnStop:  commonStop(12, nil, 0),
	}, nil
}

var _ gobs.IService = (*S12)(nil)

type S13 struct{ err error }

func (s *S13) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
		},
		OnSetup: commonSetup(13, s.err, 80),
		OnStop:  commonStop(13, nil, 0),
	}, nil
}

var _ gobs.IService = (*S13)(nil)
