package xlog

type Level int32

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

const (
	levelDebug = "debug"
	levelInfo  = "info"
	levelWarn  = "warn"
	levelError = "error"
	levelFatal = "fatal"
)

var levelToStr = map[Level]string{
	DEBUG: levelDebug,
	INFO:  levelInfo,
	WARN:  levelWarn,
	ERROR: levelError,
	FATAL: levelFatal,
}

func (l Level) String() string {
	return levelToStr[l]
}
