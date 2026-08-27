package cmd

import (
	"cmp"
	"log"
	"log/slog"

	"github.com/spf13/viper"
	"github.com/twiglab/h2o/cache"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/clog/wal"
	"github.com/twiglab/h2o/hank"
)

func rootLog() *slog.Logger {
	rlogF := viper.GetString("hank.log.root.file")
	rlogL := viper.GetString("hank.log.root.level")

	hlogL := viper.GetString("hank.log.level")
	level := clog.Level(cmp.Or(rlogL, hlogL))

	log := clog.NewLog(rlogF, level)
	slog.SetDefault(log)
	return log
}

func serverLog() *slog.Logger {
	sLogF := viper.GetString("hank.log.server.file")
	sLogL := viper.GetString("hank.log.server.level")
	hlogL := viper.GetString("hank.log.level")

	level := clog.Level(cmp.Or(sLogL, hlogL))
	l := clog.NewLog(sLogF, level)
	return l
}

func wallog() *wal.WAL {
	logf := viper.GetString("hank.wal.file")
	if logf == "" {
		log.Fatalln("wal file is null. ***MUST*** set hank.wal.file")
	}
	log.Println("wal file:", logf)
	return wal.New(wal.Conf{Filename: logf})
}

func mqtt() *hank.MQTTAction {
	broker := viper.GetString("hank.sender.mqtt.broker")
	clientID := viper.GetString("hank.sender.mqtt.client_id")

	cli, err := hank.NewMQTTClient(cmp.Or(clientID, hank.CLIENT_ID), broker)
	if err != nil {
		log.Fatal(err)
	}
	return hank.NewMQTTAction(cli)
}

func nats() *hank.NatsAction {
	url := viper.GetString("hank.sender.nats.url")
	n, err := hank.NewNatsAction(url)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func sender() hank.Sender {
	use := viper.GetString("hank.sender.use")
	switch use {
	case "mqtt":
		log.Println("using mqtt")
		return mqtt()
	case "nats":
		log.Println("using nats")
		return nats()
	}
	log.Println("using logAction")
	return hank.LogAction{}
}

func simple() hank.SimpleMD {
	proj := viper.GetString("hank.meta.simple.project")
	if proj == "" {
		log.Fatal("hank.meta.simple.project is empty")
	}
	return hank.SimpleMD{Project: proj}
}

func backend() cache.Cache[string, hank.MetaData] {
	var backend cache.Cache[string, hank.MetaData]
	b := viper.GetString("hank.meta.backend")
	switch b {
	default:
		backend = simple()
	}
	return backend
}

func enh() *hank.Enh {
	m := backend()
	return &hank.Enh{Cache: m}
}

func playback() *hank.PlayBack {
	logF := viper.GetString("hank.playback.file")
	if logF == "" {
		log.Fatalln("playback file is null. ***MUST*** set hank.playback.file")
	}
	log.Println("playback file:", logF)
	return hank.NewPlayBack(logF)
}
