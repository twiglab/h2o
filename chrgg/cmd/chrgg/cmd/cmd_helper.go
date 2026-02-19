package cmd

import (
	"log"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/chrgg"
	"github.com/twiglab/h2o/chrgg/orm"
	"github.com/twiglab/h2o/chrgg/orm/ent"
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
	log := chrgg.NewLog(logF, level)
	slog.SetDefault(log)
	return log
}

func serverLog() *slog.Logger {
	logF := viper.GetString("log.server.file")
	logL := viper.GetString("log.server.level")

	level := logLevel(logL)
	l := chrgg.NewLog(logF, level)
	return l
}

func cdrlog() *slog.Logger {
	logF := viper.GetString("chargg.cdr.file")
	if logF == "" {
		log.Fatalln("cdr file is null. ***MUST*** set chrgg.cdr.file")
	}
	log.Println("cdr file:", logF)
	return chrgg.NewLog(logF, slog.LevelInfo)
}

func mqttcli() mqtt.Client {
	broker := viper.GetString("chrgg.mqtt.broker")
	if broker == "" {
		log.Fatalf("no broker")
	}
	cli, err := chrgg.NewMQTTClient(chrgg.CLIENT_ID, broker)
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

func webaddr() string {
	addr := viper.GetString("chrgg.web.addr")
	if addr != "" {
		return addr
	}
	return ":10007"
}

func topics() map[string]byte {
	subs := viper.GetStringSlice("chrgg.mqtt.topics")
	if len(subs) == 0 {
		log.Fatal("topic is 0")
	}
	m := make(map[string]byte)
	for _, t := range subs {
		m[t] = 0x00
	}
	return m
}

func entcli() *ent.Client {
	name := viper.GetString("chrgg.db.name")
	dsn := viper.GetString("chrgg.db.dsn")

	cli, err := orm.OpenEntClient(name, dsn)
	if err != nil {
		log.Fatal(err)
	}
	return cli
}
