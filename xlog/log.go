package xlog

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

var def unsafe.Pointer

func init() {
	def = unsafe.Pointer(newLogger())
}

func GetLogger() Logger {
	return getDefLogger()
}

func getDefLogger() *logger {
	return (*logger)(atomic.LoadPointer(&def))
}

func setDefLogger(l *logger) {
	atomic.StorePointer(&def, unsafe.Pointer(l))
}

func SetLevel(level Level) {
	l := getDefLogger().clone()
	l.level = level
	setDefLogger(l)
}

func SetEnableCaller(enable bool) {
	l := getDefLogger().clone()
	l.enableCaller = enable
	setDefLogger(l)
}

func SetTimeFormat(format string) {
	l := getDefLogger().clone()
	l.timeFormat = format
	setDefLogger(l)
}

func SetEncode(e EncodeType) {
	l := getDefLogger().clone()
	if e == PLAIN {
		l.encode = plainEncode
	} else {
		l.encode = jsonEncode
	}
	setDefLogger(l)
}

func SetWriter(w Writer) {
	l := getDefLogger().clone()
	l.writer = w
	setDefLogger(l)
}

func Fatal(a ...any) {
	getDefLogger().fatalf(fmt.Sprint(a...))
}

func Fatalf(format string, a ...any) {
	getDefLogger().fatalf(fmt.Sprintf(format, a...))
}

func Fatalw(msg string, fields ...Field) {
	getDefLogger().fatalf(msg, fields...)
}

func Error(a ...any) {
	getDefLogger().error(nil, a...)
}

func Errorf(format string, a ...any) {
	getDefLogger().errorf(nil, format, a...)
}

func Errorw(msg string, fields ...Field) {
	getDefLogger().errorw(msg, fields...)
}

func Warn(a ...any) {
	getDefLogger().warn(nil, a...)
}

func Warnf(format string, a ...any) {
	getDefLogger().warnf(nil, format, a...)
}

func Warnw(msg string, fields ...Field) {
	getDefLogger().warnw(msg, fields...)
}

func Info(a ...any) {
	getDefLogger().info(nil, a...)
}

func Infof(format string, a ...any) {
	getDefLogger().infof(nil, format, a...)
}

func Infow(msg string, fields ...Field) {
	getDefLogger().infow(msg, fields...)
}

func Debug(a ...any) {
	getDefLogger().debug(nil, a...)
}

func Debugf(format string, a ...any) {
	getDefLogger().debugf(nil, format, a...)
}

func Debugw(msg string, fields ...Field) {
	getDefLogger().debugw(msg, fields...)
}
