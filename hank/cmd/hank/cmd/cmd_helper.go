package cmd

import (
	"cmp"
	"context"
	"log"
	"log/slog"

	"github.com/spf13/viper"
	"github.com/twiglab/h2o/abm"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/clog/wal"
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
	rlogF := viper.GetString("log.root.file")
	rlogL := viper.GetString("log.root.level")
	logL := viper.GetString("log.level")

	level := logLevel(cmp.Or(rlogL, logL))
	log := clog.NewLog(rlogF, level)
	slog.SetDefault(log)
	return log
}

func serverLog() *slog.Logger {
	sLogF := viper.GetString("log.server.file")
	sLogL := viper.GetString("log.server.level")
	logL := viper.GetString("log.level")

	level := logLevel(cmp.Or(sLogL, logL))
	l := clog.NewLog(sLogF, level)
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
	load := viper.GetString("hank.meta.load")
	get := viper.GetString("hank.meta.get")
	list := viper.GetString("hank.meta.list")

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

func playback() *hank.PlayBack {
	f := viper.GetString("hank.playback.file")
	logf := cmp.Or(f, "logs/playback/record.log")
	return hank.NewPlayBack(logf)
}
