package main

import (
	"context"
	"sample/services"
	"syscall"

	"github.com/xarest/gobs"
)

func main() {
	ctx := context.Background()
	bs := gobs.NewBootstrap()
	bs.AddMany(&services.A{}, &services.B{}, &services.C{})
	bs.StartBootstrap(ctx, syscall.SIGINT, syscall.SIGTERM)
}
