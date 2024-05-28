package gobs

type Config struct {
	IsConcurrent    bool
	Logger          LogFnc
	EnableLogModule bool
	EnableLogDetail bool

	EnableLogAdd   bool
	EnableLogStart bool
	EnableLogSetup bool
	EnableLogStop  bool
}

var DefaultConfig = Config{
	IsConcurrent: true,
}
