package logger

type LogFnc func(format string, args ...interface{})

type Logger struct {
	log         *LogFnc
	isLogDetail bool
	tag         string
}

func NewLog(log *LogFnc) *Logger {
	return &Logger{
		log: log,
	}
}

func (l *Logger) SetDetail(isEnabled bool) {
	l.isLogDetail = isEnabled
}

func (l *Logger) Clone() *Logger {
	return &Logger{
		log:         l.log,
		tag:         l.tag,
		isLogDetail: l.isLogDetail,
	}
}

func (l *Logger) Log(format string, args ...interface{}) {
	if l.isLogDetail {
		l.LogS(format, args...)
	}
}

func (l *Logger) LogS(format string, args ...interface{}) {
	if l.log != nil {
		args = append([]interface{}{l.tag + ":"}, args...)
		(*l.log)("%s "+format, args...)
	}
}

func (l *Logger) SetTag(tag string) func() {
	preTag := l.tag
	l.tag = tag
	return func() {
		l.tag = preTag
	}
}

func (l *Logger) AddTag(tag string) func() {
	preTag := l.tag
	l.tag = preTag + "/" + tag
	return func() {
		l.tag = preTag
	}
}
