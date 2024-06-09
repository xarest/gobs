package gobs_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/traphamxuan/gobs"
	"github.com/traphamxuan/gobs/logger"
)

func Test_Sync(t *testing.T) {
	fmt.Println("Test_Bootstrap")
	mainCtx := context.Background()
	ctx, _ := context.WithDeadline(mainCtx, time.Now().Add(5*time.Second))
	var logger logger.LogFnc = func(s string, i ...interface{}) {
		fmt.Printf(s+"\n", i...)
	}
	bs := gobs.NewBootstrap(gobs.Config{
		NumOfConcurrencies: 0,
		Logger:             &logger,
		EnableLogDetail:    true,
	})
	bs.AddDefault(&D{})

	if err := bs.Init(ctx); err != nil {
		log.Fatalf("Init expected no error, but got %v", err)
	}

	if err := bs.Setup(ctx); err != nil {
		log.Fatalf("Setup expected no error, but got %v", err)
	}

	a, ok := bs.GetService(&A{}, "").(*A)
	if !ok || a == nil {
		log.Fatal("Expected A is valid")
	}
	b, ok := bs.GetService(&B{}, "").(*B)
	if !ok || b == nil {
		log.Fatal("Expected B is valid")
	}
	c, ok := bs.GetService(&C{}, "").(*C)
	if !ok || c == nil {
		log.Fatal("Expected C is valid")
	}
	d, ok := bs.GetService(&D{}, "").(*D)
	if !ok || d == nil {
		log.Fatal("Expected D is valid")
	}

	if b.A != a {
		log.Fatalf("Expected B.A is %p, but got %p", a, b.A)
	}
	if c.A != a {
		log.Fatalf("Expected C.A is %p, but got %p", a, c.A)
	}
	if c.B != b {
		log.Fatalf("Expected C.B is %p, but got %p", b, c.B)
	}
	if d.B != b {
		log.Fatalf("Expected D.B is %p, but got %p", b, d.B)
	}
	if d.C != c {
		log.Fatalf("Expected D.C is %p, but got %p", c, d.C)
	}

	if err := bs.Start(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Fatalf("Expect context canceled, but got %v", err)
		}
	}

	if err := bs.Stop(ctx); err != nil {
		log.Fatalf("Expected no error, but got %v", err)
	}
}
