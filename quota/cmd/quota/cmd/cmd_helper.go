package cmd

import (
	"cmp"
	"log"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/clog/wal"
	"github.com/twiglab/h2o/dbcli"
	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/quota"
)

func rootLog() *slog.Logger {
	root := cmp.Or(viper.GetString("h2o.log.root"), "logs")
	level := clog.Level(viper.GetString("h2o.log.level"))

	logL := clog.NewLog("root.log", level, clog.LogRoot(root), clog.LogDir(quota.APP_ID))
	slog.SetDefault(logL)
	log.Println("root log:", root, "root.log", level)
	return logL
}

func serverLog() *slog.Logger {
	root := cmp.Or(viper.GetString("h2o.log.root"), "logs")
	level := clog.Level(viper.GetString("h2o.log.level"))

	logL := clog.NewLog("server.log", level, clog.LogRoot(root), clog.LogDir(quota.APP_ID))
	log.Println("server log:", root, "server.log", level)
	return logL
}

func cwal() *wal.WAL {
	root := cmp.Or(viper.GetString("h2o.log.wal.root"), "wal")
	log.Println("wal file:", root, "wal.log")
	return wal.New("wal.log", clog.LogRoot(root), clog.LogDir(quota.APP_ID))
}

func mqttcli() mqtt.Client {
	broker := viper.GetString("h2o.mqtt.broker")
	if broker == "" {
		log.Fatalf("no broker")
	}
	cli, err := quota.NewMQTTClient(quota.APP_ID, broker)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("mqtt - broker: %s, clientID: %s\n", broker, quota.APP_ID)
	return cli
}

func webaddr() string {
	addr := cmp.Or(viper.GetString("h2o.web.quota"), ":10003")
	log.Println("web addr:", addr)
	return addr
}

func entcli() *ent.Client {
	name := viper.GetString("h2o.db.name")
	dsn := viper.GetString("h2o.db.dsn")

	cli, err := dbcli.OpenEntClient(name, dsn)
	if err != nil {
		log.Fatal(err)
	}
	return cli
}
