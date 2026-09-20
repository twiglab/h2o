package cmd

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/nab/orm"
)

func sender(cli mqtt.Client) nab.Sender {
	use := viper.GetString("nab.sender.use")
	switch use {
	case "mqtt":
		log.Println("using mqtt")
		return nab.NewMQTTAction(cli)
	}
	log.Println("using logAction")
	return nab.LogAction{}
}

func mqttcli() mqtt.Client {
	broker := viper.GetString("nab.sender.mqtt.broker")

	if broker == "" {
		log.Fatalf("no broker")
	}

	box := viper.GetString("nab.box")
	clientID := nab.ClientID(box)
	cli, err := nab.NewMQTTClient(clientID, broker)
	if err != nil {
		log.Fatal(fmt.Errorf("mqttcli err: %w", err))
	}
	log.Printf("nab.box: %s, broker: %s, clientID: %s\n", box, broker, clientID)
	return cli
}

func global(db *orm.IDB) nab.Global {
	box := viper.GetString("nab.box")
	project := viper.GetString("nab.project")
	g := nab.Global{
		IDB:     db,
		Box:     box,
		Project: project,
	}

	if err := g.InitClients(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Printf("nab.box: %s, nab.project: %s\n", box, project)
	return g
}

func webaddr() string {
	addr := viper.GetString("nab.web.addr")
	return cmp.Or(addr, ":10008")
}

func db() *orm.IDB {
	dev := viper.GetString("nab.idb.dev")
	cli := viper.GetString("nab.idb.cli")

	db, err := orm.NewIDB(dev, cli)
	if err != nil {
		log.Fatal(fmt.Errorf("idb err: %w", err))
	}
	return db
}

func rootLog() *slog.Logger {
	rlogF := viper.GetString("nab.log.root.file")
	rlogL := viper.GetString("nab.log.root.level")
	logL := viper.GetString("nab.log.level")

	level := clog.Level(cmp.Or(rlogL, logL))
	log := clog.NewLog(rlogF, level)
	slog.SetDefault(log)
	return log
}

func serverLog() *slog.Logger {
	sLogF := viper.GetString("nab.log.server.file")
	sLogL := viper.GetString("nab.log.server.level")
	logL := viper.GetString("nab.log.level")

	level := clog.Level(cmp.Or(sLogL, logL))
	l := clog.NewLog(sLogF, level)
	return l
}
