package gobs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/traphamxuan/gobs"
)

func Test_Bootstrap(t *testing.T) {
	ctx := context.TODO()
	bs := gobs.NewBootstrap(gobs.Config{
		IsConcurrent: true,
		Logger: func(s string, i ...interface{}) {
			fmt.Printf(s+"\n", i...)
		},
		EnableLogModule: true,
		EnableLogDetail: true,
		EnableLogAdd:    true,
		EnableLogSetup:  true,
		EnableLogStart:  true,
		EnableLogStop:   true,
	})
	bs.AddDefault(&D{})

	if err := bs.Setup(ctx); err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	a, ok := bs.GetService(&A{}, "").(*A)
	if !ok || a == nil {
		t.Error("Expected A is valid")
	}
	b, ok := bs.GetService(&B{}, "").(*B)
	if !ok || b == nil {
		t.Error("Expected B is valid")
	}
	c, ok := bs.GetService(&C{}, "").(*C)
	if !ok || c == nil {
		t.Error("Expected C is valid")
	}
	d, ok := bs.GetService(&D{}, "").(*D)
	if !ok || d == nil {
		t.Error("Expected D is valid")
	}

	if b.A != a {
		t.Errorf("Expected B.A is %p, but got %p", a, b.A)
	}
	if c.A != a {
		t.Errorf("Expected C.A is %p, but got %p", a, c.A)
	}
	if c.B != b {
		t.Errorf("Expected C.B is %p, but got %p", b, c.B)
	}
	if d.B != b {
		t.Errorf("Expected D.B is %p, but got %p", b, d.B)
	}
	if d.C != c {
		t.Errorf("Expected D.C is %p, but got %p", c, d.C)
	}
}
