package xlog

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// def stores a pointer to the default logger instance.
var def unsafe.Pointer

// init sets the default logger to a new logger instance.
func init() {
	def = unsafe.Pointer(newLogger())
}

// getDefLogger returns a pointer to the default logger instance.
func getDefLogger() *logger {
	return (*logger)(atomic.LoadPointer(&def))
}

// setDefLogger sets the default logger instance to the given logger.
func setDefLogger(l *logger) {
	atomic.StorePointer(&def, unsafe.Pointer(l))
}

// GetLogger returns the default logger instance.
func GetLogger() Logger {
	return getDefLogger()
}

// SetLevel sets the logging level for the default logger.
func SetLevel(level Level) {
	l := getDefLogger().clone()
	l.level = level
	setDefLogger(l)
}

// SetCaller sets whether or not to log the caller's function name and line number for the default logger.
func SetCaller(enable bool) {
	l := getDefLogger().clone()
	l.enableCaller = enable
	setDefLogger(l)
}

// SetColor sets whether or not to use colorized output levelTag on plain encoding for the default logger.
func SetColor(enable bool) {
	l := getDefLogger().clone()
	l.enableColor = enable
	setDefLogger(l)
}

// SetTimeFormat sets the time format string for the default logger.
func SetTimeFormat(format string) {
	l := getDefLogger().clone()
	l.timeFormat = format
	setDefLogger(l)
}

// SetEncode sets the log encoding type for the default logger.
func SetEncode(e EncodeType) {
	l := getDefLogger().clone()
	if e == PLAIN {
		l.encode = plainEncode
	} else {
		l.encode = jsonEncode
	}
	setDefLogger(l)
}

// SetWriter sets the log writer for the default logger.
func SetWriter(w Writer) {
	l := getDefLogger().clone()
	l.writer = w
	setDefLogger(l)
}

// Log writes a log entry with the given content and options to the default logger.
// The content is of any type and can be any loggable object, such as a string or a struct.
// The options are optional and can be used to customize the log entry.
// This function delegates to the Log method of the default logger instance.
func Log(content any, opts ...LogOption) {
	opts = append(opts, WithCallerSkipOne)
	getDefLogger().Log(content, opts...)
}

// Fatal logs a message at fatal level and exits the program with an error status.
func Fatal(a ...any) {
	getDefLogger().fatalf(fmt.Sprint(a...))
}

// Fatalf logs a formatted message at fatal level and exits the program with an error status.
func Fatalf(format string, a ...any) {
	getDefLogger().fatalf(fmt.Sprintf(format, a...))
}

// Fatalw logs a message with extra fields at fatal level and exits the program with an error status.
func Fatalw(msg string, fields ...Field) {
	getDefLogger().fatalf(msg, fields...)
}

// Error logs a message at error level.
func Error(a ...any) {
	getDefLogger().error(nil, a...)
}

// Errorf logs a formatted message at error level.
func Errorf(format string, a ...any) {
	getDefLogger().errorf(nil, format, a...)
}

// Errorw logs a message with extra fields at error level.
func Errorw(msg string, fields ...Field) {
	getDefLogger().errorw(msg, fields...)
}

// Warn logs a message at warning level.
func Warn(a ...any) {
	getDefLogger().warn(nil, a...)
}

// Warnf logs a formatted message at warning level.
func Warnf(format string, a ...any) {
	getDefLogger().warnf(nil, format, a...)
}

// Warnw logs a message with extra fields at warning level.
func Warnw(msg string, fields ...Field) {
	getDefLogger().warnw(msg, fields...)
}

// Info logs a message at info level.
func Info(a ...any) {
	getDefLogger().info(nil, a...)
}

// Infof logs a formatted message at info level.
func Infof(format string, a ...any) {
	getDefLogger().infof(nil, format, a...)
}

// Infow logs a message with extra fields at info level.
func Infow(msg string, fields ...Field) {
	getDefLogger().infow(msg, fields...)
}

// Debug logs a message at debug level.
func Debug(a ...any) {
	getDefLogger().debug(nil, a...)
}

// Debugf logs a formatted message at debug level.
func Debugf(format string, a ...any) {
	getDefLogger().debugf(nil, format, a...)
}

// Debugw logs a message with extra fields at debug level.
func Debugw(msg string, fields ...Field) {
	getDefLogger().debugw(msg, fields...)
}
