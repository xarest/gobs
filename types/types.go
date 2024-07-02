package types

import (
	"context"

	"github.com/traphamxuan/gobs/common"
)

type ITask interface {
	Run(ctx context.Context, status common.ServiceStatus) error
	IsRunAsync(status common.ServiceStatus) bool
	Name() string
	DependOn(status common.ServiceStatus) []ITask
	Followers(status common.ServiceStatus) []ITask
}
