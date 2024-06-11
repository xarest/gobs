package gobs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/traphamxuan/gobs"
)

//	var l logger.LogFnc = func(s string, i ...interface{}) {
//		fmt.Printf(s+"\n", i...)
//	}
// var setupOrder []int

func Test_AsyncPerformance(t *testing.T) {
	numOfDependencies := 20
	level := 5
	// shared := 5
	ctx := context.Background()
	service := NewSampleAsyncService(numOfDependencies, level)
	bs := gobs.NewBootstrap(gobs.Config{
		NumOfConcurrencies: -1,
		// Logger:             &l,
		// EnableLogDetail:    true,
	})
	bs.AddDefault(service)
	if err := bs.Init(ctx); err != nil {
		t.Fatalf("Error: %v", err)
	}
	if err := bs.Setup(ctx); err != nil {
		t.Fatalf("Error: %v", err)
	}
	// if len(setupOrder) != numOfServices {
	// 	t.Fatalf("Expected %d services to be setup, got %d", numOfServices, len(setupOrder))
	// }
	// if err := bs.Start(ctx); err != nil {
	// 	t.Fatalf("Error: %v", err)
	// }
	// if err := bs.Stop(ctx); err != nil {
	// 	t.Fatalf("Error: %v", err)
	// }
	fmt.Printf("Test_Performance Async with %d services\n", numOfServices)
}

type SampleAsyncService struct {
	level     int
	id        int
	numOfDeps int
	// dependencies []gobs.IService
}

// Init implements gobs.IService.
func (s *SampleAsyncService) Init(ctx context.Context, c *gobs.Component) error {
	numOfServices++
	if s.level > 0 {
		deps := make([]SampleAsyncService, s.numOfDeps)
		c.ExtraDeps = make([]gobs.CustomService, s.numOfDeps)
		newLevel := s.level - 1
		for i := 0; i < s.numOfDeps; i++ {
			deps[i].level = newLevel
			deps[i].numOfDeps = s.numOfDeps
			deps[i].id = numOfServices
			c.ExtraDeps[i] = gobs.CustomService{
				Name:     fmt.Sprintf("Sample-%d-%d-%d", newLevel, i, numOfServices),
				Instance: &deps[i],
			}
		}
	}
	onSetup := func(context.Context, []gobs.IService, []gobs.CustomService) error {
		return nil
	}
	c.OnSetupAsync = &onSetup

	return nil
}

var _ gobs.IService = (*SampleAsyncService)(nil)

func NewSampleAsyncService(numOfDeps, level int) gobs.IService {
	return &SampleAsyncService{
		level:     level,
		id:        0,
		numOfDeps: numOfDeps,
	}
}
