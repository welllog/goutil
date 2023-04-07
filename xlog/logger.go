package xlog

import (
	"fmt"
	"os"
	"time"
)

// logger represents a logger instance with configurable options
type logger struct {
	level        Level                        // the minimum level of logging to output
	enableCaller bool                         // flag indicating whether to log the caller information
	enableColor  bool                         // flag indicating whether to use colorized output for levelTag on plain encoding
	timeFormat   string                       // time format to use for logging
	encode       func(o *logOption, w Writer) // encoding function to use for logging
	writer       Writer                       // writer to output log to
}

// NewLogger returns a new Logger instance with optional configurations
func NewLogger(opts ...LoggerOption) Logger {
	return newLogger(opts...)
}

// newLogger creates a new logger instance with default options that can be customized with the provided options
func newLogger(opts ...LoggerOption) *logger {
	l := logger{
		level:        DEBUG,
		enableCaller: true,
		enableColor:  true,
		timeFormat:   time.RFC3339,
		encode:       jsonEncode,
		writer:       csWriter,
	}
	for _, opt := range opts {
		opt(&l)
	}
	return &l
}

// LoggerOption is a functional option type for configuring a logger instance
type LoggerOption func(*logger)

// WithLoggerLevel sets the minimum logging level for the logger instance
func WithLoggerLevel(level Level) LoggerOption {
	return func(l *logger) {
		l.level = level
	}
}

// WithLoggerCaller sets whether to log the caller information in the output
func WithLoggerCaller(enable bool) LoggerOption {
	return func(l *logger) {
		l.enableCaller = enable
	}
}

// WithLoggerColor sets whether to use colorized output for levelTag on plain encoding
func WithLoggerColor(enable bool) LoggerOption {
	return func(l *logger) {
		l.enableColor = enable
	}
}

// WithLoggerTimeFormat sets the time format to use for logging
func WithLoggerTimeFormat(format string) LoggerOption {
	return func(l *logger) {
		l.timeFormat = format
	}
}

// WithLoggerEncode sets the encoding type to use for logging
func WithLoggerEncode(e EncodeType) LoggerOption {
	return func(l *logger) {
		if e == PLAIN {
			l.encode = plainEncode
		} else {
			l.encode = jsonEncode
		}
	}
}

// WithLoggerWriter sets the output writer for the logger instance
func WithLoggerWriter(w Writer) LoggerOption {
	return func(l *logger) {
		l.writer = w
	}
}

func (l *logger) Log(content any, opts ...LogOption) {
	o := logOption{
		enableCaller: l.enableCaller,
		callerSkip:   defCallerSkip - 1,
	}
	for _, opt := range opts {
		opt(&o)
	}
	o.content = content
	if l.IsEnabled(o.level) {
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
	if l.IsEnabled(ERROR) {
		l.output(&logOption{
			level:        ERROR,
			enableCaller: l.enableCaller,
			content:      fmt.Sprint(a...),
			fields:       fields,
		})
	}
}

func (l *logger) errorf(fields []Field, format string, a ...any) {
	if l.IsEnabled(ERROR) {
		l.output(&logOption{
			level:        ERROR,
			enableCaller: l.enableCaller,
			content:      fmt.Sprintf(format, a...),
			fields:       fields,
		})
	}
}

func (l *logger) errorw(v any, fields ...Field) {
	if l.IsEnabled(ERROR) {
		l.output(&logOption{
			level:        ERROR,
			enableCaller: l.enableCaller,
			content:      v,
			fields:       fields,
		})
	}
}

func (l *logger) warn(fields []Field, a ...any) {
	if l.IsEnabled(WARN) {
		l.output(&logOption{
			level:        WARN,
			enableCaller: l.enableCaller,
			content:      fmt.Sprint(a...),
			fields:       fields,
		})
	}
}

func (l *logger) warnf(fields []Field, format string, a ...any) {
	if l.IsEnabled(WARN) {
		l.output(&logOption{
			level:        WARN,
			enableCaller: l.enableCaller,
			content:      fmt.Sprintf(format, a...),
			fields:       fields,
		})
	}
}

func (l *logger) warnw(v any, fields ...Field) {
	if l.IsEnabled(WARN) {
		l.output(&logOption{
			level:        WARN,
			enableCaller: l.enableCaller,
			content:      v,
			fields:       fields,
		})
	}
}

func (l *logger) info(fields []Field, a ...any) {
	if l.IsEnabled(INFO) {
		l.output(&logOption{
			level:        INFO,
			enableCaller: l.enableCaller,
			content:      fmt.Sprint(a...),
			fields:       fields,
		})
	}
}

func (l *logger) infof(fields []Field, format string, a ...any) {
	if l.IsEnabled(INFO) {
		l.output(&logOption{
			level:        INFO,
			enableCaller: l.enableCaller,
			content:      fmt.Sprintf(format, a...),
			fields:       fields,
		})
	}
}

func (l *logger) infow(v any, fields ...Field) {
	if l.IsEnabled(INFO) {
		l.output(&logOption{
			level:        INFO,
			enableCaller: l.enableCaller,
			content:      v,
			fields:       fields,
		})
	}
}

func (l *logger) debug(fields []Field, a ...any) {
	if l.IsEnabled(DEBUG) {
		l.output(&logOption{
			level:        DEBUG,
			enableCaller: l.enableCaller,
			content:      fmt.Sprint(a...),
			fields:       fields,
		})
	}
}

func (l *logger) debugf(fields []Field, format string, a ...any) {
	if l.IsEnabled(DEBUG) {
		l.output(&logOption{
			level:        DEBUG,
			enableCaller: l.enableCaller,
			content:      fmt.Sprintf(format, a...),
			fields:       fields,
		})
	}
}

func (l *logger) debugw(v any, fields ...Field) {
	if l.IsEnabled(DEBUG) {
		l.output(&logOption{
			level:        DEBUG,
			enableCaller: l.enableCaller,
			content:      v,
			fields:       fields,
		})
	}
}

func (l *logger) IsEnabled(level Level) bool {
	return level >= l.level
}

func (l *logger) output(o *logOption) {
	if l.enableCaller {
		if o.callerSkip <= 0 {
			o.callerSkip = defCallerSkip
		}
		o.caller = getCaller(o.callerSkip)
	}
	if o.levelTag == "" {
		o.levelTag = o.level.String()
	}
	o.enableColor = l.enableColor
	o.timestamp = time.Now().Format(l.timeFormat)
	l.encode(o, l.writer)
}

func (l *logger) clone() *logger {
	return &logger{
		level:        l.level,
		enableCaller: l.enableCaller,
		enableColor:  l.enableColor,
		timeFormat:   l.timeFormat,
		encode:       l.encode,
		writer:       l.writer,
	}
}
