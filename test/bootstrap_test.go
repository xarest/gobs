package gobs_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/traphamxuan/gobs"
)

type BootstrapSuit struct {
	suite.Suite
}

func TestBootstrap(t *testing.T) {
	suite.Run(t, new(BootstrapSuit))
}

// func (s *SchedulerSuit) SetupSuite() {
// 	fmt.Println("SchedulerSuit/SetupSuite")
// }

// func (s *SchedulerSuit) TearDownSuite() {
// 	fmt.Println("SchedulerSuit/TearDownSuite")
// }

// func (s *SchedulerSuit) SetupTest() {
// 	fmt.Println("SchedulerSuit/SetupTest")
// }

// func (s *SchedulerSuit) TearDownTest() {
// 	fmt.Println("SchedulerSuit/TearDownTest")
// }

func (s *BootstrapSuit) TestSync() {
	t := s.T()
	mainCtx := context.Background()
	ctx, cancel := context.WithDeadline(mainCtx, time.Now().Add(5*time.Second))
	defer cancel() // It's a good practice to call cancel even if not strictly necessary here

	// var logger logger.LogFnc = func(s string, i ...interface{}) {
	// 	fmt.Printf(s+"\n", i...)
	// }
	bs := gobs.NewBootstrap(gobs.Config{
		NumOfConcurrencies: 0,
		Logger:             &log,
		// EnableLogDetail: true,
	})
	bs.AddDefault(&D{})

	require.NoError(t, bs.Init(ctx), "Init expected no error")
	require.NoError(t, bs.Setup(ctx), "Setup expected no error")

	a, ok := bs.GetService(&A{}, "").(*A)
	require.True(t, ok, "Expected A is valid")
	require.NotNil(t, a, "Expected A is not nil")

	b, ok := bs.GetService(&B{}, "").(*B)
	require.True(t, ok, "Expected B is valid")
	require.NotNil(t, b, "Expected B is not nil")

	c, ok := bs.GetService(&C{}, "").(*C)
	require.True(t, ok, "Expected C is valid")
	require.NotNil(t, c, "Expected C is not nil")

	d, ok := bs.GetService(&D{}, "").(*D)
	require.True(t, ok, "Expected D is valid")
	require.NotNil(t, d, "Expected D is not nil")

	assert.Equal(t, a, b.A, "Expected B.A is equal to A")
	assert.Equal(t, a, c.A, "Expected C.A is equal to A")
	assert.Equal(t, b, c.B, "Expected C.B is equal to B")
	assert.Equal(t, b, d.B, "Expected D.B is equal to B")
	assert.Equal(t, c, d.C, "Expected D.C is equal to C")

	require.NoError(t, bs.Start(ctx), "Expected no error from Start")
	require.NoError(t, bs.Stop(ctx), "Expected no error from Stop")
}
