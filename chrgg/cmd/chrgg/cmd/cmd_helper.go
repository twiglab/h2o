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
)

func rootLog() *slog.Logger {
	rlogF := viper.GetString("chrgg.log.root.file")
	rlogL := viper.GetString("chrgg.log.root.level")
	logL := viper.GetString("chrgg.log.level")

	level := clog.Level(cmp.Or(rlogL, logL))
	log := clog.NewLog(rlogF, level)
	slog.SetDefault(log)
	return log
}

func serverLog() *slog.Logger {
	sLogF := viper.GetString("chrgg.log.server.file")
	sLogL := viper.GetString("chrgg.log.server.level")
	logL := viper.GetString("chrgg.log.level")

	level := clog.Level(cmp.Or(sLogL, logL))
	l := clog.NewLog(sLogF, level)
	return l
}

func cwal() *wal.WAL {
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
	return cmp.Or(addr, ":10003")
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
