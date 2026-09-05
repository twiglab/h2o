package cmd

import (
	"cmp"
	"log"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/chrgg"
	"github.com/twiglab/h2o/chrgg/orm"
	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/clog/wal"
	"github.com/twiglab/h2o/pkg/common"
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
		common.WaterTopic:       0x01,
		common.ElectricityTopic: 0x01,
		common.GasTopic:         0x01,
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

func dbx() *orm.DBx {
	return &orm.DBx{Cli: entcli()}
}

func cs() *chrgg.ChargeServer {
	return &chrgg.ChargeServer{
		DBx: dbx(),

		Logger: serverLog(),
	}
}
