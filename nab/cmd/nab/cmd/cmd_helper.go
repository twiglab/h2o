package cmd

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"path"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/nab/orm"
	"github.com/twiglab/h2o/nab/orm/ent"
)

func mqttLog() {
	root := cmp.Or(viper.GetString("nab.log.root"), "logs")
	level := clog.Level(viper.GetString("nab.log.level"))

	logF := path.Join(root, "mqtt.log")
	out := clog.NewLogWriter(logF)
	l := log.New(out, "mqtt - ", log.LstdFlags)
	switch level {
	case slog.LevelDebug:
		mqtt.DEBUG = l
		mqtt.WARN = l
		mqtt.ERROR = l
		mqtt.CRITICAL = l
	case slog.LevelInfo:
		mqtt.WARN = l
		mqtt.ERROR = l
		mqtt.CRITICAL = l
	case slog.LevelWarn:
		mqtt.WARN = l
		mqtt.ERROR = l
		mqtt.CRITICAL = l
	case slog.LevelError:
		mqtt.ERROR = l
		mqtt.CRITICAL = l
	}
}

func mqttcli(clientID, topic string, h mqtt.MessageHandler) (mqtt.Client, nab.Sender) {
	broker := viper.GetString("nab.sender.mqtt.broker")

	if broker == "" {
		log.Fatalf("no broker")
	}

	opts := mqtt.NewClientOptions()
	opts.SetClientID(clientID)
	opts.AddBroker(broker)

	// 开关指令允许丢失且无顺序要求（README），故关闭 paho 的报文顺序保证；
	// 注意 order=false 时每条消息独立 goroutine，handler 内不得阻塞。
	opts.SetOrderMatters(false)

	// 允许丢失消息，每次链接都用新session
	opts.SetCleanSession(true)

	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetConnectTimeout(10 * time.Second)
	opts.SetMaxReconnectInterval(2 * time.Minute)

	opts.SetWriteTimeout(10 * time.Second)

	opts.SetOnConnectHandler(func(mc mqtt.Client) {
		slog.Debug("mqtt - onConnect")
		token := mc.Subscribe(topic, 0x00, h)
		if err := mqtt.WaitTokenTimeout(token, 5*time.Second); err != nil {
			slog.Error("mqtt 订阅失败", slog.Any("error", err))
			return
		}
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		slog.Error("mqtt - Connect lost", slog.Any("error", err))
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

func onoffLog() *slog.Logger {
	root := cmp.Or(viper.GetString("nab.log.root"), "logs")
	logL := clog.NewLog("onoff.log", slog.LevelInfo, clog.LogRoot(root))
	log.Println("onoff log:", root, "onoff.log", slog.LevelInfo)
	return logL
}

func pool() *ants.Pool {
	p, err := ants.NewPool(32, ants.WithNonblocking(true))
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func randText() string {
	now := time.Now()
	return now.Format(f)
}

const f = "20060102150405"

func httpserver(h http.Handler) *http.Server {
	srv := &http.Server{
		Addr:    webaddr(),
		Handler: h,
	}
	return srv
}
