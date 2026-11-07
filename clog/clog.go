package clog

import (
	"io"
	"log/slog"
	"os"
	"path"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	WARN  = "warn"
)

const CONSOLE = "console"

func Level(s string) slog.Level {
	switch strings.ToLower(s) {
	case DEBUG:
		return slog.LevelDebug
	case INFO:
		return slog.LevelInfo
	case ERROR:
		return slog.LevelError
	case WARN:
		return slog.LevelWarn
	}
	return slog.LevelInfo
}

func isConsole(logFile string) bool {
	if logFile == "" || logFile == CONSOLE {
		return true
	}
	return false
}

func ConsoleLog(level slog.Level) *slog.Logger {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	l := slog.New(h)
	slog.SetDefault(l)
	return l
}

func NewLogWriter(logf string) io.Writer {
	return &lumberjack.Logger{
		Filename:   logf,
		MaxSize:    30,  // megabytes
		MaxAge:     60,  // days
		MaxBackups: 100, // 最多100个
		LocalTime:  true,
	}
}

type LogOpts struct {
	Root string
	Dir  string
}

type LogOpt func(o *LogOpts)

func LogRoot(root string) LogOpt {
	return func(o *LogOpts) {
		o.Root = root
	}
}

func LogDir(dir string) LogOpt {
	return func(o *LogOpts) {
		o.Dir = dir
	}
}

func NewLog(logFile string, level slog.Level, opts ...LogOpt) *slog.Logger {
	lo := &LogOpts{}
	for _, f := range opts {
		f(lo)
	}

	var out io.Writer = os.Stdout
	if !isConsole(logFile) {
		logF := path.Join(lo.Root, lo.Dir, logFile)
		out = NewLogWriter(logF)
	}
	h := slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level})
	return slog.New(h)
}
