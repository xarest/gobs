package gobs

import "github.com/traphamxuan/gobs/logger"

type Config struct {
	NumOfConcurrencies int
	Logger             *logger.LogFnc
	EnableLogDetail    bool
}

var DefaultConfig = Config{
	NumOfConcurrencies: 0,
}
