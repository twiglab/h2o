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
	logF := viper.GetString("log.root.file")
	logL := viper.GetString("log.root.level")

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
	logF := viper.GetString("datalog.file")
	if logF == "" {
		logF = "logs/datalog.log"
	}
	log.Println("datalog file:", logF)
	return hank.NewLog(logF, slog.LevelInfo)
}

func sender() hank.Sender {
	broker := viper.GetString("mqtt.broker")
	if broker == "" {
		log.Println("using logAction ...")
		return hank.LogAction{}
	}
	log.Println("using mqttAction ... broker", broker)
	cli, err := hank.NewMQTTClient("hank-plugin", broker)
	if err != nil {
		log.Fatal(err)
	}

	return hank.NewMQTTAction(cli)
}
