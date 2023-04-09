//go:build windows

package xlog

// writeLevelWithColor is a stub function that always returns the input levelTag string unchanged. It's used as a placeholder
// function on Windows OS because the ANSI color codes that are used to format console output on Unix-based systems aren't
// supported by the Windows console.
func writeLevelWithColor(level Level, levelTag string, buf *Buffer) {
	_, _ = buf.WriteString(levelTag)
}
