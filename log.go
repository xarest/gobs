package gobs

type LogFnc func(format string, args ...interface{})

func (bs *Bootstrap) LogModule(isLog bool, message string, args ...interface{}) {
	if bs.config.Logger != nil && isLog && bs.config.EnableLogModule {
		bs.config.Logger(message, args...)
	}
}

func (sb *Component) LogComponent(message string, args ...interface{}) {
	if sb.config.Logger != nil && sb.config.EnableLogDetail {
		sb.config.Logger(message, args...)
	}
}
