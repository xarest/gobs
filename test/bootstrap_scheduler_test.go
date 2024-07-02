package gobs_test

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traphamxuan/gobs"
	"github.com/traphamxuan/gobs/logger"
)

var log logger.LogFnc = func(s string, i ...interface{}) {
	fmt.Printf(s+"\n", i...)
}

func (s *BootstrapSuit) TestSyncScheduler() {
	t := s.T()
	setupOrder = []int{}
	bs := gobs.NewBootstrap(gobs.Config{
		NumOfConcurrencies: 0,
	})
	ctx := context.TODO()
	require.NoError(t, bs.AddDefault(new(S1)), "AddDefault expected no error")
	require.NoError(t, bs.Init(ctx), "Init expected no error")
	require.NoError(t, bs.Setup(ctx), "Setup expected no error")

	expectedBootOrder := []int{9, 10, 4, 11, 5, 2, 6, 12, 7, 13, 8, 3, 1}
	require.Equal(t, len(expectedBootOrder), len(setupOrder), "Expected setupOrder length to match expectedBootOrder length")
	assert.Equal(t, expectedBootOrder, setupOrder, "Expected setupOrder to match expectedBootOrder")
}

func (s *BootstrapSuit) TestAsyncScheduler() {
	t := s.T()
	setupOrder = []int{}
	bs := gobs.NewBootstrap(gobs.Config{
		NumOfConcurrencies: gobs.DEFAULT_MAX_CONCURRENT,
	})
	ctx := context.TODO()
	require.NoError(t, bs.AddDefault(new(S1)), "AddDefault expected no error")
	require.NoError(t, bs.Init(ctx), "Init expected no error")
	require.NoError(t, bs.Setup(ctx), "Setup expected no error")

	expectedBootOrder := []int{11, 12, 7, 13, 8, 10, 6, 3, 9, 4, 5, 2, 1}
	require.Equal(t, len(expectedBootOrder), len(setupOrder), "Expected setupOrder length to match expectedBootOrder length")
	assert.Equal(t, expectedBootOrder, setupOrder, "Expected setupOrder to match expectedBootOrder")
}

func (s *BootstrapSuit) TestAsyncSchedulerWithError() {
	t := s.T()
	setupOrder = []int{}
	bs := gobs.NewBootstrap(gobs.Config{
		NumOfConcurrencies: gobs.DEFAULT_MAX_CONCURRENT,
		// Logger:             &log,
		// EnableLogDetail:    true,
	})
	ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(5*time.Second))
	defer cancel()
	setupOrder = []int{}
	require.NoError(t, bs.AddDefault(new(S1)), "AddDefault expected no error")
	s9 := &S9{err: assert.AnError}
	require.NoError(t, bs.AddDefault(s9), "AddDefault expected no error")
	require.NoError(t, bs.Init(ctx), "Init expected no error")
	s9, ok := bs.GetService(&S9{}, "").(*S9)
	require.True(t, ok, "Expected GetService return S9")
	require.NotNil(t, s9, "Expected S9 is not nil")
	s9.err = assert.AnError
	require.Error(t, bs.Setup(ctx), "Setup expected error")
	expectedBootOrder := []int{11, 12, 7, 13, 8, 10, 6, 3, 9}
	require.Equal(t, len(expectedBootOrder), len(setupOrder), "Expected setupOrder length to match expectedBootOrder length")
	assert.Equal(t, expectedBootOrder, setupOrder, "Expected setupOrder to match expectedBootOrder")

	require.NoError(t, bs.Stop(ctx), "Setup expected no error")

	expectedStopOrder := []int{3, 6, 7, 8, 10, 11, 12, 13}
	require.Equal(t, len(expectedStopOrder), len(stopOrder), "Expected stopOrder length to match expectedStopOrder length")
	assert.Equal(t, expectedStopOrder, stopOrder, "Expected stopOrder to match expectedStopOrder")
	s9.err = nil
}
