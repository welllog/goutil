//go:build windows

package xlog

func wrapLevelWithColor(level Level, levelTag string) string {
	return levelTag
}
