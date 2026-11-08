package cmd

import (
	"cmp"
	"log"
	"log/slog"

	"github.com/spf13/viper"
	"github.com/twiglab/h2o/archon/orm"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/dbcli"
	"github.com/twiglab/h2o/dbcli/ent"
)

func rootLog() *slog.Logger {
	root := cmp.Or(viper.GetString("h2o.log.root"), "logs")
	level := clog.Level(viper.GetString("h2o.log.level"))

	logL := clog.NewLog("root.log", level, clog.LogRoot(root), clog.LogDir("archon"))
	slog.SetDefault(logL)
	log.Println("root log:", root, "root.log", level)
	return logL
}

func serverLog() *slog.Logger {
	root := cmp.Or(viper.GetString("h2o.log.root"), "logs")
	level := clog.Level(viper.GetString("h2o.log.level"))

	logL := clog.NewLog("server.log", level, clog.LogRoot(root), clog.LogDir("archon"))
	log.Println("server log:", root, "server.log", level)
	return logL
}

func webaddr() string {
	addr := cmp.Or(viper.GetString("h2o.web.quota"), ":10008")
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

func dbx(c *ent.Client) *orm.DBx {
	return &orm.DBx{Client: c}
}
