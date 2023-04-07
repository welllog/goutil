package xlog

// Logger is an interface that defines the methods for logging.
type Logger interface {
	// Log writes a log message with the given content and options.
	Log(content any, opts ...LogOption)

	// Fatal writes a log message with the FATAL log level and call os.Exit(1).
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Fatalw(msg string, fields ...Field)

	// Error writes a log message with the ERROR log level.
	Error(args ...any)
	Errorf(format string, args ...any)
	Errorw(msg string, fields ...Field)

	// Warn writes a log message with the WARN log level.
	Warn(args ...any)
	Warnf(format string, args ...any)
	Warnw(msg string, fields ...Field)

	// Info writes a log message with the INFO log level.
	Info(args ...any)
	Infof(format string, args ...any)
	Infow(msg string, fields ...Field)

	// Debug writes a log message with the DEBUG log level.
	Debug(args ...any)
	Debugf(format string, args ...any)
	Debugw(msg string, fields ...Field)

	// IsEnabled returns whether the given log level is enabled or not.
	IsEnabled(level Level) bool
}

// Field is a struct that represents a key-value pair of additional data to include in a log message.
type Field struct {
	Key   string
	Value any
}

// LogOption is a function that modifies a logOption struct.
type LogOption func(*logOption)

// logOption is a struct that represents options to use when logging a message.
type logOption struct {
	level        Level // level is the severity level of the log message.
	enableCaller bool  // enableCaller indicates whether to include caller information in the log message.
	enableColor  bool  // enableColor indicates whether to enable colorized output for the levelTag on plain encoding.
	callerSkip   int   // callerSkip is the number of stack frames to skip to find the caller information.

	// levelTag is the string representation of the severity level
	// The default debug, info, warn, error, and fatal correspond to DEBUG, INFO, WARN, ERROR, and FATAL log levels respectively
	// users can also customize semantic tags, such as slow.
	levelTag  string
	timestamp string  // timestamp is the time the log message was created.
	caller    string  // caller is the file and line number where the log message was created.
	content   any     // content is the main content of the log message.
	fields    []Field // fields is a slice of key-value pairs of additional data to include in the log message.
}

// defCallerSkip is the default number of stack frames to skip to find the caller information.
const defCallerSkip = 4

// WithLevel returns a LogOption that sets the logging level and the corresponding tag.
func WithLevel(level Level, levelTag string) LogOption {
	return func(o *logOption) {
		o.level = level
		o.levelTag = levelTag
	}
}

// WithCaller returns a LogOption that enables or disables logging the caller information.
func WithCaller(enable bool) LogOption {
	return func(o *logOption) {
		o.enableCaller = enable
	}
}

// WithCallerSkip returns a LogOption that sets the number of stack frames to skip when logging caller information.
func WithCallerSkip(skip int) LogOption {
	return func(o *logOption) {
		o.callerSkip += skip
	}
}

// WithFields returns a LogOption that appends additional key-value pairs to the logged message.
func WithFields(fields ...Field) LogOption {
	return func(o *logOption) {
		o.fields = append(o.fields, fields...)
	}
}

// WithCallerSkipOne is a LogOption that increments the number of stack frames to skip by 1 when logging caller information.
func WithCallerSkipOne(o *logOption) {
	o.callerSkip++
}

// WithCallerSkipTwo is a LogOption that increments the number of stack frames to skip by 2 when logging caller information.
func WithCallerSkipTwo(o *logOption) {
	o.callerSkip += 2
}

func withFatal(o *logOption) {
	o.level = FATAL
}

func withError(o *logOption) {
	o.level = ERROR
}

func withWarn(o *logOption) {
	o.level = WARN
}

func withInfo(o *logOption) {
	o.level = INFO
}

func withDebug(o *logOption) {
	o.level = DEBUG
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
