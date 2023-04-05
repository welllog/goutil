//go:build !windows

package xlog

const (
	Reset     = "\033[0m"
	Red       = "\033[31m"
	RedBold   = "\033[1;31m"
	Green     = "\033[32m"
	GreenBold = "\033[1;32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	BlueBold  = "\033[1;34m"
	Purple    = "\033[35m"
	Cyan      = "\033[36m"
	Gray      = "\033[37m"
	White     = "\033[97m"
)

func wrapLevelWithColor(level Level, levelTag string) string {
	switch level {
	case FATAL:
		return RedBold + levelTag + Reset
	case ERROR:
		return Red + levelTag + Reset
	case WARN:
		return Yellow + levelTag + Reset
	case INFO:
		return Green + levelTag + Reset
	case DEBUG:
		return Gray + levelTag + Reset
	}
	return levelTag
}
