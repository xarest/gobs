package gobs

import "github.com/traphamxuan/gobs/logger"

type Config struct {
	IsConcurrent    bool
	Logger          *logger.LogFnc
	EnableLogDetail bool
}

var DefaultConfig = Config{
	IsConcurrent: true,
}
