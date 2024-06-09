package scheduler_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/traphamxuan/gobs"
	"github.com/traphamxuan/gobs/logger"
)

func Test_SyncScheduler(t *testing.T) {
	setupOrder = []int{}
	fmt.Println("Start test Test_SyncScheduler")
	var logger logger.LogFnc = func(s string, i ...interface{}) {
		fmt.Printf(s+"\n", i...)
	}
	bs := gobs.NewBootstrap(gobs.Config{
		NumOfConcurrencies: 0,
		Logger:             &logger,
		EnableLogDetail:    true,
	})
	ctx := context.Background()
	bs.AddDefault(new(S1))
	bs.Init(ctx)
	bs.Setup(ctx)
	expectedBootOrder := []int{9, 10, 4, 5, 2, 11, 6, 12, 7, 13, 8, 3, 1}
	if len(setupOrder) != len(expectedBootOrder) {
		t.Fatalf("Expected %d, but got %d", len(expectedBootOrder), len(setupOrder))
	}
	for i, orderId := range setupOrder {
		if orderId != expectedBootOrder[i] {
			t.Errorf("Expected %d, but got %d", expectedBootOrder[i], orderId)
		}
	}
}
