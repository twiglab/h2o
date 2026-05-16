package cmd

import (
	"cmp"
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/spf13/viper"
	"github.com/twiglab/h2o/abm"
	"github.com/twiglab/h2o/cache"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/clog/wal"
	"github.com/twiglab/h2o/hank"
	"github.com/twiglab/h2o/hank/hkv"
)

func rootLog() *slog.Logger {
	rlogF := viper.GetString("log.root.file")
	rlogL := viper.GetString("log.root.level")
	logL := viper.GetString("log.level")

	level := clog.Level(cmp.Or(rlogL, logL))
	log := clog.NewLog(rlogF, level)
	slog.SetDefault(log)
	return log
}

func serverLog() *slog.Logger {
	sLogF := viper.GetString("log.server.file")
	sLogL := viper.GetString("log.server.level")
	logL := viper.GetString("log.level")

	level := clog.Level(cmp.Or(sLogL, logL))
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
	load := viper.GetString("hank.meta.ddb.load")
	get := viper.GetString("hank.meta.ddb.get")
	list := viper.GetString("hank.meta.ddb.list")

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

func bhkv() cache.Cache[string, hank.MetaData] {
	log.Println("backend:", hkv.Key)

	dbname := viper.GetString("hank.meta.hkv.dbname")
	dsn := viper.GetString("hank.meta.hkv.dsn")
	sqlget := viper.GetString("hank.meta.hkv.sql_get")
	project := viper.GetString("hank.meta.hkv.project")

	logf := viper.GetString("hank.meta.hkv.logfile")
	logl := viper.GetString("hank.meta.hkv.loglevel")

	conf := hkv.HankDBConf{
		DBName:  dbname,
		DSN:     dsn,
		Project: project,
		SQLGet:  sqlget,
		Logger:  clog.NewLog(logf, clog.Level(logl)),
	}

	hdb, err := hkv.NewHankDB(conf)
	if err != nil {
		log.Fatal(err)
	}
	return cache.WithCache(hdb, hkv.NewCache(3*time.Hour))
}

func backend() cache.Cache[string, hank.MetaData] {
	var backend cache.Cache[string, hank.MetaData]
	b := viper.GetString("hank.meta.backend")
	switch b {
	case hkv.Key:
		backend = bhkv()
	case "ddb":
		backend, _ = ddb()
	default:
		backend = cache.EmptyCache[string, hank.MetaData]{}
	}
	return backend
}

func enh() *hank.Enh {
	m := backend()
	return &hank.Enh{Cache: m}
}

func playback() *hank.PlayBack {
	f := viper.GetString("hank.playback.file")
	logf := cmp.Or(f, "logs/playback/record.log")
	return hank.NewPlayBack(logf)
}
