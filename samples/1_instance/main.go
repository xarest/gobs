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
	bs.AddDefault(&services.Service{})
	bs.StartBootstrap(ctx, syscall.SIGINT, syscall.SIGTERM)
}
