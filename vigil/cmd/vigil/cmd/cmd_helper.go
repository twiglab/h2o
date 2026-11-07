package cmd

import (
	"cmp"
	"fmt"
	"log"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/clog/wal"
	"github.com/twiglab/h2o/dbcli"
	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/vigil"
	"github.com/twiglab/h2o/vigil/orm"
)

func mqttcli() mqtt.Client {
	broker := viper.GetString("h2o.mqtt.broker")
	if broker == "" {
		log.Fatalf("no broker")
	}
	cli, err := vigil.NewMQTTClient(vigil.APP_ID, broker)
	if err != nil {
		log.Fatal(fmt.Errorf("mqttcli err: %w", err))
	}
	return cli
}

func webaddr() string {
	addr := viper.GetString("h2o.web.vigil")
	return cmp.Or(addr, ":10007")
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

func dbx(c *ent.Client) *orm.DBx {
	return &orm.DBx{Client: c}
}

/*
	func tdb() *tsdb.Schemaless {
		dsn := viper.GetString("vigil.tsdb.dsn")
		sch, err := tsdb.NewSchLe(dsn)
		if err != nil {
			log.Fatal(fmt.Errorf("SchLe err: %w", err))
		}
		return sch
	}
*/

func cwal() *wal.WAL {
	root := cmp.Or(viper.GetString("h2o.log.wal.root"), "wal")
	log.Println("wal file:", root, "wal.log")
	return wal.New("wal.log", clog.LogRoot(root), clog.LogDir(vigil.APP_ID))
}

func rootLog() *slog.Logger {
	root := cmp.Or(viper.GetString("h2o.log.root"), "logs")
	level := clog.Level(viper.GetString("h2o.log.level"))

	logL := clog.NewLog("root.log", level, clog.LogRoot(root), clog.LogDir(vigil.APP_ID))
	slog.SetDefault(logL)
	log.Println("root log:", root, "root.log", level)
	return logL
}

func serverLog() *slog.Logger {
	root := cmp.Or(viper.GetString("h2o.log.root"), "logs")
	level := clog.Level(viper.GetString("h2o.log.level"))

	logL := clog.NewLog("server.log", level, clog.LogRoot(root), clog.LogDir(vigil.APP_ID))
	log.Println("server log:", root, "server.log", level)
	return logL
}
