package cmd

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"log/slog"
	"path"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/clog/wal"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/nab/orm"
	"github.com/twiglab/h2o/nab/orm/ent"
)

func mqttcli(topic string, h mqtt.MessageHandler) (mqtt.Client, nab.Sender) {
	broker := viper.GetString("nab.sender.mqtt.broker")

	if broker == "" {
		log.Fatalf("no broker")
	}

	box := viper.GetString("nab.box")
	clientID := nab.ClientID(box)

	opts := mqtt.NewClientOptions()
	opts.SetClientID(clientID)
	opts.AddBroker(broker)
	opts.SetOrderMatters(false)

	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetConnectTimeout(10 * time.Second)
	opts.SetMaxReconnectInterval(2 * time.Minute)

	opts.SetWriteTimeout(10 * time.Second)

	opts.SetOnConnectHandler(func(mc mqtt.Client) {
		slog.Debug("mqtt - onConnect")
		token := mc.Subscribe(topic, 0x00, h)
		if token.WaitTimeout(5*time.Second) && token.Error() != nil {
			slog.Error("mqtt - 订阅失败", slog.Any("error", token.Error()))
		}
	})

	client := mqtt.NewClient(opts)
	return client, nab.NewMQTTAction(client)
}

func global() nab.Global {
	box := viper.GetString("nab.box")
	if box == "" {
		log.Fatal("nab.box is nil")
	}

	project := viper.GetString("nab.project")
	if project == "" {
		log.Fatal("nab.project is nil")
	}

	g := nab.Global{
		Box:     box,
		Project: project,
	}

	log.Printf("global - nab.box: %s, nab.project: %s\n", box, project)
	return g
}

func clientMgr(db *orm.IDB) *nab.ClientMgr {
	m, err := nab.NewClientMgr(db)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func webaddr() string {
	addr := viper.GetString("nab.web.addr")
	return cmp.Or(addr, ":10008")
}

func dbfile(g nab.Global) (df, cf string) {
	p := cmp.Or(viper.GetString("nab.idb.dir"), "idb")
	dev := cmp.Or(viper.GetString("nab.idb.dev"), g.Box+"_dev.csv")
	cli := cmp.Or(viper.GetString("nab.idb.cli"), g.Box+"_cli.csv")
	return path.Join(p, dev), path.Join(p, cli)
}

func db(g nab.Global) *orm.IDB {
	dev, cli := dbfile(g)

	db, err := orm.NewIDB(dev, cli)
	if err != nil {
		log.Fatal(fmt.Errorf("idb err: %w", err))
	}
	log.Printf("idb - nab.idb.dev: %s, nab.idb.cli: %s\n", dev, cli)
	return db
}

func taskInterval() time.Duration {
	i := cmp.Or(viper.GetDuration("nab.loop.task.interval"), 1000) * time.Millisecond
	log.Println("nab.loop.task.interval:", i)
	return i
}

func loops(db *orm.IDB, f func([]*ent.Dev) nab.Job) nab.Job {
	delay := cmp.Or(viper.GetDuration("nab.loops.delay"), 1000) * time.Millisecond

	log.Println("nab.loops.delay:", delay)

	lps := nab.NewLoops()

	devs, err := db.AllDev(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	a := nab.SplitByCli(devs)

	for _, s := range a {
		lps.AddToNewLoop(delay, f(s))
	}

	return lps
}

func rootLog() *slog.Logger {
	root := cmp.Or(viper.GetString("nab.log.root"), "logs")
	level := clog.Level(viper.GetString("nab.log.level"))

	logL := clog.NewLog("root.log", level, clog.LogRoot(root))
	slog.SetDefault(logL)

	log.Println("root log:", root, "root.log", level)
	return logL
}

func serverLog() *slog.Logger {
	root := cmp.Or(viper.GetString("nab.log.root"), "logs")
	level := clog.Level(viper.GetString("nab.log.level"))

	logL := clog.NewLog("server.log", level, clog.LogRoot(root))

	log.Println("server log:", root, "server.log", level)
	return logL
}

func cwal() *wal.WAL {
	root := cmp.Or(viper.GetString("nab.log.wal.root"), "wal")
	log.Println("wal file:", root, "wal.log")
	return wal.New("wal.log", clog.LogRoot(root))
}

func pool() *ants.Pool {
	p, err := ants.NewPool(10)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
