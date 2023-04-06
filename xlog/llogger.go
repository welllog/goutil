package xlog

type Logger interface {
	Log(content any, opts ...LogOption)

	Fatal(a ...any)
	Fatalf(format string, a ...any)
	Fatalw(msg string, fields ...Field)

	Error(a ...any)
	Errorf(format string, a ...any)
	Errorw(msg string, fields ...Field)

	Warn(a ...any)
	Warnf(format string, a ...any)
	Warnw(msg string, fields ...Field)

	Info(a ...any)
	Infof(format string, a ...any)
	Infow(msg string, fields ...Field)

	Debug(a ...any)
	Debugf(format string, a ...any)
	Debugw(msg string, fields ...Field)

	IsOut(level Level) bool
}

type Field struct {
	Key   string
	Value any
}

type LogOption func(*logOption)

type logOption struct {
	level        Level
	levelTag     string
	enableCaller bool
	callerSkip   int
	time         string
	caller       string
	content      any
	fields       []Field
}

const callerSkip = 4

func WithLevel(level Level, levelTag string) LogOption {
	return func(o *logOption) {
		o.level = level
		o.levelTag = levelTag
	}
}

func WithCaller(enable bool) LogOption {
	return func(o *logOption) {
		o.enableCaller = enable
	}
}

func WithCallerSkip(skip int) LogOption {
	return func(o *logOption) {
		o.callerSkip += skip
	}
}

func WithFields(fields ...Field) LogOption {
	return func(o *logOption) {
		o.fields = append(o.fields, fields...)
	}
}

func WithCallerSkipOne(o *logOption) {
	o.callerSkip++
}

func WithCallerSkipTwo(o *logOption) {
	o.callerSkip += 2
}

var (
	fieldTime    = "@timestamp"
	fieldLevel   = "level"
	fieldContent = "content"
	fieldCaller  = "caller"

	filterField = map[string]struct{}{
		fieldTime:    {},
		fieldLevel:   {},
		fieldContent: {},
		fieldCaller:  {},
	}
)
