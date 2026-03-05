package cmd

import (
	"cmp"
	"context"
	"log"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/abm"
	"github.com/twiglab/h2o/chrgg"
	"github.com/twiglab/h2o/chrgg/orm"
	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/pkg/common"
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

func cdrWal() *wal.WAL {
	logF := viper.GetString("chrgg.wal.file")
	if logF == "" {
		log.Fatalln("cdr file is null. ***MUST*** set chrgg.wal.file")
	}
	log.Println("wal file:", logF)
	return wal.New(wal.Conf{Filename: logF})
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
	return cmp.Or(addr, ":10007")
}

func topics() map[string]byte {
	return map[string]byte{
		common.WaterTopic:       0x00,
		common.ElectricityTopic: 0x00,
		common.GasTopic:         0x00,
	}
}

func entcli() *ent.Client {
	name := viper.GetString("chrgg.db.name")
	dsn := viper.GetString("chrgg.db.dsn")

	//cli, err := orm.OpenEntClient(name, dsn, ent.Debug())
	cli, err := orm.OpenEntClient(name, dsn)
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

func dbx() *chrgg.DBx {
	return &chrgg.DBx{Cli: entcli()}
}

func ddb() (*abm.DuckABM[string, chrgg.AloneRuler], abm.Conf) {
	load := viper.GetString("chrgg.ce.alone.load")
	get := viper.GetString("chrgg.ce.alone.get")
	list := viper.GetString("chrgg.ce.alone.list")

	c := abm.Conf{
		LoadSQL: load,
		GetSQL:  get,
		ListSQL: list,
		Period:  60,
	}

	db, err := abm.NewDuckABM[string, chrgg.AloneRuler](c)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Loop(context.Background()); err != nil {
		log.Fatal(err)
	}
	return db, c
}

func ce() chrgg.ChargeEngine {
	b := viper.GetString("chrgg.ce.backend")
	switch b {
	case "alone":
		log.Println("ce: alone")
		feedb, _ := ddb()
		return chrgg.NewAloneEngine(feedb)
	}
	log.Println("ce: EngZ")
	return chrgg.EngZ
}

func cs() *chrgg.ChargeServer {
	return &chrgg.ChargeServer{
		CdrWAL:      cdrWal(),
		DBx:         dbx(),
		ChargEngine: ce(),
		CheckFunc:   chrgg.DefaultCheck,
		SkipFunc:    chrgg.DefaultSkip,

		Logger: serverLog(),
	}
}
