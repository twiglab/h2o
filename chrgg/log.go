package chrgg

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func isConsole(logFile string) bool {
	if logFile == "" || logFile == "console" {
		return true
	}
	return false
}

func NewLog(logFile string, level slog.Level) *slog.Logger {
	var out io.Writer = os.Stdout
	if !isConsole(logFile) {
		out = &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    10, // megabytes
			MaxBackups: 10,
			MaxAge:     10, //days
		}
	}
	h := slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level})
	return slog.New(h)
}
