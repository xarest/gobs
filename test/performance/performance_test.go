package performance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/traphamxuan/gobs"
)

var numOfServices = 0

// var logger gobs.LogFnc = func(s string, i ...interface{}) {
// 	fmt.Printf(s+"\n", i...)
// }

func Test_SyncPerformance(t *testing.T) {
	numOfDependencies := 20
	level := 5
	// shared := 5
	ctx := context.Background()
	service := NewSampleService(numOfDependencies, level)
	bs := gobs.NewBootstrap(gobs.Config{
		NumOfConcurrencies: 0,
		// Logger:             &logger,
		// EnableLogSchedule:  true,
		// EnableLogDetail:    true,
	})
	bs.AddDefault(service)
	if err := bs.Init(ctx); err != nil {
		t.Fatalf("Error: %v", err)
	}
	if err := bs.Setup(ctx); err != nil {
		t.Fatalf("Error: %v", err)
	}
	// if err := bs.Start(ctx); err != nil {
	// 	t.Fatalf("Error: %v", err)
	// }
	// if err := bs.Stop(ctx); err != nil {
	// 	t.Fatalf("Error: %v", err)
	// }
	fmt.Printf("Test_Performance with %d services\n", numOfServices)
}

type SampleService struct {
	level     int
	id        int
	numOfDeps int
	// dependencies []gobs.IService
}

// Init implements gobs.IService.
func (s *SampleService) Init(ctx context.Context, c *gobs.Component) error {
	if s.level == 0 {
		return nil
	}
	deps := make([]SampleService, s.numOfDeps)
	c.ExtraDeps = make([]gobs.CustomService, s.numOfDeps)
	newLevel := s.level - 1
	for i := 0; i < s.numOfDeps; i++ {
		numOfServices++
		deps[i].level = newLevel
		deps[i].numOfDeps = s.numOfDeps
		deps[i].id = numOfServices
		c.ExtraDeps[i] = gobs.CustomService{
			Name:     fmt.Sprintf("Sample-%d-%d", newLevel, numOfServices),
			Instance: &deps[i],
		}
	}
	onSetup := func(context.Context, []gobs.IService, []gobs.CustomService) error {
		return nil
	}
	c.OnSetup = &onSetup

	return nil
}

var _ gobs.IService = (*SampleService)(nil)

func NewSampleService(numOfDeps, level int) gobs.IService {
	return &SampleService{
		level:     level,
		id:        0,
		numOfDeps: numOfDeps,
	}
}
