package cmd

import (
	"log"
	"log/slog"

	"github.com/spf13/viper"
	"github.com/twiglab/h2o/hank"
)

func logLevel(s string) slog.Level {
	switch s {
	case "debug", "DEBUG":
		return slog.LevelDebug
	case "info", "INFO":
		return slog.LevelInfo
	case "error", "ERROR":
		return slog.LevelError
	case "warn", "WARN":
		return slog.LevelWarn
	}
	return slog.LevelInfo
}

func rootLog() *slog.Logger {
	logF := viper.GetString("log.file")
	logL := viper.GetString("log.level")

	level := logLevel(logL)
	log := hank.NewLog(logF, level)
	slog.SetDefault(log)
	return log
}

func serverLog() *slog.Logger {
	logF := viper.GetString("log.server.file")
	logL := viper.GetString("log.server.level")

	level := logLevel(logL)
	l := hank.NewLog(logF, level)
	return l
}

func dataLog() *slog.Logger {
	return hank.NewLog("logs/datalog.log", slog.LevelInfo)
}

func sender() hank.Sender {
	borker := viper.GetString("mqtt.borker")
	if borker == "" {
		log.Println("using logAction")
		return hank.NewLogAction()
	}
	cli, err := hank.NewMQTTClient("hank-plugin", borker)
	if err != nil {
		log.Fatal(err)
	}

	return hank.NewMQTTAction(cli)
}
