package gobs_test

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traphamxuan/gobs"
)

// var logger logger.LogFnc = func(s string, i ...interface{}) {
// 	fmt.Printf(s+"\n", i...)
// }

func (s *BootstrapSuit) TestSyncScheduler() {
	t := s.T()
	setupOrder = []int{}
	bs := gobs.NewBootstrap(gobs.Config{
		IsConcurrent: false,
	})
	ctx := context.TODO()
	require.NoError(t, bs.AddDefault(new(S1)), "AddDefault expected no error")
	require.NoError(t, bs.Init(ctx), "Init expected no error")
	require.NoError(t, bs.Setup(ctx), "Setup expected no error")

	expectedBootOrder := []int{9, 10, 4, 11, 5, 2, 6, 12, 7, 13, 8, 3, 1}
	require.Equal(t, len(expectedBootOrder), len(setupOrder), "Expected setupOrder length to match expectedBootOrder length")
	assert.Equal(t, expectedBootOrder, setupOrder, "Expected setupOrder to match expectedBootOrder")
}
