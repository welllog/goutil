package xlog

import (
	"fmt"
	"os"
	"time"
)

type EncodeType int

const (
	JSON EncodeType = iota
	PLAIN
)

type logger struct {
	level        Level
	enableCaller bool
	timeFormat   string
	encode       func(o *logOption, w Writer)
	writer       Writer
}

func NewLogger(opts ...LoggerOption) Logger {
	return newLogger(opts...)
}

func newLogger(opts ...LoggerOption) *logger {
	l := logger{
		level:        DEBUG,
		enableCaller: true,
		timeFormat:   time.RFC3339,
		encode:       jsonEncode,
		writer:       csWriter,
	}
	for _, opt := range opts {
		opt(&l)
	}
	return &l
}

type LoggerOption func(*logger)

func WithLoggerLevel(level Level) LoggerOption {
	return func(l *logger) {
		l.level = level
	}
}

func WithLoggerCaller(enable bool) LoggerOption {
	return func(l *logger) {
		l.enableCaller = enable
	}
}

func WithLoggerTimeFormat(format string) LoggerOption {
	return func(l *logger) {
		l.timeFormat = format
	}
}

func WithLoggerEncode(e EncodeType) LoggerOption {
	return func(l *logger) {
		if e == PLAIN {
			l.encode = plainEncode
		} else {
			l.encode = jsonEncode
		}
	}
}

func WithLoggerWriter(w Writer) LoggerOption {
	return func(l *logger) {
		l.writer = w
	}
}

func (l *logger) Log(content any, opts ...LogOption) {
	o := logOption{
		enableCaller: l.enableCaller,
		callerSkip:   callerSkip - 1,
	}
	for _, opt := range opts {
		opt(&o)
	}
	o.content = content
	if l.IsOut(o.level) {
		l.output(&o)
	}
	if o.level == FATAL {
		os.Exit(1)
	}
}

func (l *logger) Fatal(a ...any) {
	l.fatalf(fmt.Sprint(a...))
}

func (l *logger) Fatalf(format string, a ...any) {
	l.fatalf(fmt.Sprintf(format, a...))
}

func (l *logger) Fatalw(msg string, fields ...Field) {
	l.fatalf(msg, fields...)
}

func (l *logger) Error(a ...any) {
	l.error(nil, a...)
}

func (l *logger) Errorf(format string, a ...any) {
	l.errorf(nil, format, a...)
}

func (l *logger) Errorw(msg string, fields ...Field) {
	l.errorw(msg, fields...)
}

func (l *logger) Warn(a ...any) {
	l.warn(nil, a...)
}

func (l *logger) Warnf(format string, a ...any) {
	l.warnf(nil, format, a...)
}

func (l *logger) Warnw(msg string, fields ...Field) {
	l.warnw(msg, fields...)
}

func (l *logger) Info(a ...any) {
	l.info(nil, a...)
}

func (l *logger) Infof(format string, a ...any) {
	l.infof(nil, format, a...)
}

func (l *logger) Infow(msg string, fields ...Field) {
	l.infow(msg, fields...)
}

func (l *logger) Debug(a ...any) {
	l.debug(nil, a...)
}

func (l *logger) Debugf(format string, a ...any) {
	l.debugf(nil, format, a...)
}

func (l *logger) Debugw(msg string, fields ...Field) {
	l.debugw(msg, fields...)
}

func (l *logger) fatalf(msg string, fields ...Field) {
	l.output(&logOption{
		level:        FATAL,
		enableCaller: l.enableCaller,
		content:      msg,
		fields:       fields,
	})
	os.Exit(1)
}

func (l *logger) error(fields []Field, a ...any) {
	if l.IsOut(ERROR) {
		l.output(&logOption{
			level:        ERROR,
			enableCaller: l.enableCaller,
			content:      fmt.Sprint(a...),
			fields:       fields,
		})
	}
}

func (l *logger) errorf(fields []Field, format string, a ...any) {
	if l.IsOut(ERROR) {
		l.output(&logOption{
			level:        ERROR,
			enableCaller: l.enableCaller,
			content:      fmt.Sprintf(format, a...),
			fields:       fields,
		})
	}
}

func (l *logger) errorw(v any, fields ...Field) {
	if l.IsOut(ERROR) {
		l.output(&logOption{
			level:        ERROR,
			enableCaller: l.enableCaller,
			content:      v,
			fields:       fields,
		})
	}
}

func (l *logger) warn(fields []Field, a ...any) {
	if l.IsOut(WARN) {
		l.output(&logOption{
			level:        WARN,
			enableCaller: l.enableCaller,
			content:      fmt.Sprint(a...),
			fields:       fields,
		})
	}
}

func (l *logger) warnf(fields []Field, format string, a ...any) {
	if l.IsOut(WARN) {
		l.output(&logOption{
			level:        WARN,
			enableCaller: l.enableCaller,
			content:      fmt.Sprintf(format, a...),
			fields:       fields,
		})
	}
}

func (l *logger) warnw(v any, fields ...Field) {
	if l.IsOut(WARN) {
		l.output(&logOption{
			level:        WARN,
			enableCaller: l.enableCaller,
			content:      v,
			fields:       fields,
		})
	}
}

func (l *logger) info(fields []Field, a ...any) {
	if l.IsOut(INFO) {
		l.output(&logOption{
			level:        INFO,
			enableCaller: l.enableCaller,
			content:      fmt.Sprint(a...),
			fields:       fields,
		})
	}
}

func (l *logger) infof(fields []Field, format string, a ...any) {
	if l.IsOut(INFO) {
		l.output(&logOption{
			level:        INFO,
			enableCaller: l.enableCaller,
			content:      fmt.Sprintf(format, a...),
			fields:       fields,
		})
	}
}

func (l *logger) infow(v any, fields ...Field) {
	if l.IsOut(INFO) {
		l.output(&logOption{
			level:        INFO,
			enableCaller: l.enableCaller,
			content:      v,
			fields:       fields,
		})
	}
}

func (l *logger) debug(fields []Field, a ...any) {
	if l.IsOut(DEBUG) {
		l.output(&logOption{
			level:        DEBUG,
			enableCaller: l.enableCaller,
			content:      fmt.Sprint(a...),
			fields:       fields,
		})
	}
}

func (l *logger) debugf(fields []Field, format string, a ...any) {
	if l.IsOut(DEBUG) {
		l.output(&logOption{
			level:        DEBUG,
			enableCaller: l.enableCaller,
			content:      fmt.Sprintf(format, a...),
			fields:       fields,
		})
	}
}

func (l *logger) debugw(v any, fields ...Field) {
	if l.IsOut(DEBUG) {
		l.output(&logOption{
			level:        DEBUG,
			enableCaller: l.enableCaller,
			content:      v,
			fields:       fields,
		})
	}
}

func (l *logger) IsOut(level Level) bool {
	return level >= l.level
}

func (l *logger) output(o *logOption) {
	if l.enableCaller {
		if o.callerSkip <= 0 {
			o.callerSkip = callerSkip
		}
		o.caller = getCaller(o.callerSkip)
	}
	if o.levelTag == "" {
		o.levelTag = o.level.String()
	}
	o.time = time.Now().Format(l.timeFormat)
	l.encode(o, l.writer)
}

func (l *logger) clone() *logger {
	return &logger{
		level:        l.level,
		enableCaller: l.enableCaller,
		timeFormat:   l.timeFormat,
		encode:       l.encode,
		writer:       l.writer,
	}
}
