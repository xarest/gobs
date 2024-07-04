package gobs_test

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/require"
	"github.com/traphamxuan/gobs"
)

var numOfServices = 0

func (s *BootstrapSuit) TestSyncPerformance() {
	t := s.T()
	numOfDependencies := 10
	level := 3
	ctx := context.Background()
	service := NewSampleService(numOfDependencies, level)
	bs := gobs.NewBootstrap(gobs.Config{
		NumOfConcurrencies: 0,
	})
	bs.AddDefault(service)
	require.NoError(t, bs.Init(ctx), "Bootstrap initialization failed")
	require.NoError(t, bs.Setup(ctx), "Bootstrap Setup failed")
	require.NoError(t, bs.Start(ctx), "Bootstrap Start failed")
	require.NoError(t, bs.Stop(ctx), "Bootstrap Stop failed")
}

type SampleService struct {
	level     int
	id        int
	numOfDeps int
	// dependencies []gobs.IService
}

// Init implements gobs.IService.
func (s *SampleService) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	if s.level == 0 {
		return nil, nil
	}
	sCfg := &gobs.ServiceLifeCycle{}
	deps := make([]SampleService, s.numOfDeps)
	sCfg.ExtraDeps = make([]gobs.CustomService, s.numOfDeps)
	newLevel := s.level - 1
	for i := 0; i < s.numOfDeps; i++ {
		numOfServices++
		deps[i].level = newLevel
		deps[i].numOfDeps = s.numOfDeps
		deps[i].id = numOfServices
		sCfg.ExtraDeps[i] = gobs.CustomService{
			Name:     fmt.Sprintf("Sample-%d-%d", newLevel, numOfServices),
			Instance: &deps[i],
		}
	}

	return sCfg, nil
}

var _ gobs.IService = (*SampleService)(nil)

func NewSampleService(numOfDeps, level int) gobs.IService {
	return &SampleService{
		level:     level,
		id:        0,
		numOfDeps: numOfDeps,
	}
}
