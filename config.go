package gobs

import "github.com/xarest/gobs/logger"

type Config struct {
	NumOfConcurrencies int
	Logger             logger.LogFnc
	EnableLogDetail    bool
}

const DEFAULT_MAX_CONCURRENT = -1

var DefaultConfig = Config{
	NumOfConcurrencies: DEFAULT_MAX_CONCURRENT,
}
