package gobs_test

import (
	"context"
	"fmt"
	"time"

	"github.com/xarest/gobs"
	"github.com/xarest/gobs/common"
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
		Deps: gobs.Dependencies{new(S2), new(S3)},
	}, nil
}

func (s *S1) Setup(ctx context.Context, deps ...gobs.IService) error {
	if err := gobs.Dependencies(deps).Assign(&s.S2, &s.S3); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	setupOrder = append(setupOrder, 1)
	return s.err
}

func (s *S1) Stop(ctx context.Context) error {
	return commonStop(1, nil, 0)(ctx)
}

var _ gobs.IServiceInit = (*S1)(nil)

type S2 struct {
	err error
	S4  *S4
	S5  *S5
}

func (s *S2) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S4), new(S5)},
	}, nil
}

func (s *S2) Setup(ctx context.Context, deps ...gobs.IService) error {
	if err := gobs.Dependencies(deps).Assign(&s.S4, &s.S5); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	setupOrder = append(setupOrder, 2)
	return s.err
}

func (s *S2) Stop(ctx context.Context) error {
	return commonStop(2, nil, 0)(ctx)
}

var _ gobs.IServiceInit = (*S2)(nil)

type S3 struct {
	err error
	S6  *S6
	S7  *S7
	S8  *S8
}

func (s *S3) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S6), new(S7), new(S8)},
	}, nil
}

func (s *S3) Setup(ctx context.Context, deps ...gobs.IService) error {
	if err := gobs.Dependencies(deps).Assign(&s.S6, &s.S7, &s.S8); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	setupOrder = append(setupOrder, 3)
	return s.err
}

func (s *S3) Stop(ctx context.Context) error {
	return commonStop(3, nil, 0)(ctx)
}

var _ gobs.IServiceInit = (*S3)(nil)

type S4 struct {
	err error
	S9  *S9
	S10 *S10
}

func (s *S4) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S9), new(S10)},
	}, nil
}

func (s *S4) Setup(ctx context.Context, deps ...gobs.IService) error {
	if err := gobs.Dependencies(deps).Assign(&s.S9, &s.S10); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	setupOrder = append(setupOrder, 4)
	return s.err
}

func (s *S4) Stop(ctx context.Context) error {
	return commonStop(4, nil, 0)(ctx)
}

var _ gobs.IServiceInit = (*S4)(nil)

type S5 struct {
	err error
	S9  *S9
	S10 *S10
	S11 *S11
}

func (s *S5) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S9), new(S10), new(S11)},
	}, nil
}

func (s *S5) Setup(ctx context.Context, deps ...gobs.IService) error {
	if err := gobs.Dependencies(deps).Assign(&s.S9, &s.S10, &s.S11); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	setupOrder = append(setupOrder, 5)
	return s.err
}

func (s *S5) Stop(ctx context.Context) error {
	return commonStop(5, nil, 0)(ctx)
}

var _ gobs.IServiceInit = (*S5)(nil)

type S6 struct {
	err error
	S10 *S10
	S11 *S11
}

func (s *S6) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S10), new(S11)},
	}, nil
}

func (s *S6) Setup(ctx context.Context, deps ...gobs.IService) error {
	if err := gobs.Dependencies(deps).Assign(&s.S10, &s.S11); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	setupOrder = append(setupOrder, 6)
	return s.err
}

func (s *S6) Stop(ctx context.Context) error {
	return commonStop(6, nil, 0)(ctx)
}

var _ gobs.IServiceInit = (*S6)(nil)

type S7 struct {
	err error
	S12 *S12
}

func (s *S7) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S12)},
	}, nil
}

func (s *S7) Setup(ctx context.Context, deps ...gobs.IService) error {
	if err := gobs.Dependencies(deps).Assign(&s.S12); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	setupOrder = append(setupOrder, 7)
	return s.err
}

func (s *S7) Stop(ctx context.Context) error {
	return commonStop(7, nil, 0)(ctx)
}

var _ gobs.IServiceInit = (*S7)(nil)

type S8 struct {
	err error
	S13 *S13
}

func (s *S8) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{new(S13)},
	}, nil
}

func (s *S8) Setup(ctx context.Context, deps ...gobs.IService) error {
	if err := gobs.Dependencies(deps).Assign(&s.S13); err != nil {
		fmt.Println("Failed to assign dependencies", err)
		return err
	}
	setupOrder = append(setupOrder, 8)
	return s.err
}

func (s *S8) Stop(ctx context.Context) error {
	return commonStop(8, nil, 0)(ctx)
}

var _ gobs.IServiceInit = (*S8)(nil)

type S9 struct{ err error }

func (s *S9) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
		},
	}, nil
}

func (s *S9) Setup(ctx context.Context, deps ...gobs.IService) error {
	return commonSetup(9, s.err, 100)(ctx, deps)
}

func (s *S9) Stop(ctx context.Context) error {
	return commonStop(9, nil, 0)(ctx)
}

var _ gobs.IServiceInit = (*S9)(nil)

type S10 struct{ err error }

func (s *S10) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
		},
	}, nil
}

func (s *S10) Setup(ctx context.Context, deps ...gobs.IService) error {
	return commonSetup(10, s.err, 90)(ctx, deps)
}

func (s *S10) Stop(ctx context.Context) error {
	return commonStop(10, nil, 0)(ctx)
}

var _ gobs.IServiceInit = (*S10)(nil)

type S11 struct{ err error }

func (s *S11) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
		},
	}, nil
}

func (s *S11) Setup(ctx context.Context, deps ...gobs.IService) error {
	return commonSetup(11, s.err, 1)(ctx, deps)
}

func (s *S11) Stop(ctx context.Context) error {
	return commonStop(11, nil, 0)(ctx)
}

var _ gobs.IServiceInit = (*S11)(nil)

type S12 struct{ err error }

func (s *S12) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
		},
	}, nil
}

func (s *S12) Setup(ctx context.Context, deps ...gobs.IService) error {
	return commonSetup(12, s.err, 2)(ctx, deps)
}

func (s *S12) Stop(ctx context.Context) error {
	return commonStop(12, nil, 0)(ctx)
}

var _ gobs.IServiceInit = (*S12)(nil)

type S13 struct{ err error }

func (s *S13) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		AsyncMode: map[common.ServiceStatus]bool{
			common.StatusSetup: true,
		},
	}, nil
}

func (s *S13) Setup(ctx context.Context, deps ...gobs.IService) error {
	return commonSetup(13, s.err, 80)(ctx, deps)
}

func (s *S13) Stop(ctx context.Context) error {
	return commonStop(13, nil, 0)(ctx)
}

var _ gobs.IServiceInit = (*S13)(nil)
