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
	bs.AddOrPanic(&services.S1{})
	bs.StartBootstrap(ctx, syscall.SIGINT, syscall.SIGTERM)
}
