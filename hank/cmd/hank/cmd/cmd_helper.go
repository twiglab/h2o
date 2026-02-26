package cmd

import (
	"context"
	"log"
	"log/slog"

	"github.com/spf13/viper"
	"github.com/twiglab/h2o/abm"
	"github.com/twiglab/h2o/hank"
	"github.com/twiglab/h2o/wal"
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

func wallog() *wal.WAL {
	logf := viper.GetString("hank.wal.file")
	if logf == "" {
		log.Fatalln("wal file is null. ***MUST*** set datalog.file")
	}
	log.Println("wal file:", logf)
	return wal.New(wal.Conf{Filename: logf})
}

func sender() hank.Sender {
	broker := viper.GetString("hank.mqtt.broker")
	if broker == "" {
		log.Println("using logAction ...")
		return hank.LogAction{}
	}
	log.Println("using mqttAction ... broker", broker)
	cli, err := hank.NewMQTTClient(hank.CLIENT_ID, broker)
	if err != nil {
		log.Fatal(err)
	}

	return hank.NewMQTTAction(cli)
}

func ddb() (*abm.DuckABM[string, hank.MetaData], abm.Conf) {
	load := viper.GetString("hank.abm.load")
	get := viper.GetString("hank.abm.get")
	list := viper.GetString("hank.abm.list")

	c := abm.Conf{
		LoadSQL: load,
		GetSQL:  get,
		ListSQL: list,
		Period:  60,
	}

	db, err := abm.NewDuckABM[string, hank.MetaData](c)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Loop(context.Background()); err != nil {
		log.Fatal(err)
	}
	return db, c
}

func enh() *hank.Enh {
	m, _ := ddb()
	return &hank.Enh{DDB: m}
}
