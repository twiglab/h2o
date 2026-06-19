package cmd

import (
	"cmp"
	"log"
	"log/slog"

	"github.com/spf13/viper"
	"github.com/twiglab/h2o/archon/orm"
	"github.com/twiglab/h2o/archon/orm/ent"
	"github.com/twiglab/h2o/clog"
)

func webaddr() string {
	addr := viper.GetString("archon.web.addr")
	return cmp.Or(addr, ":10008")
}

func entcli() *ent.Client {
	name := viper.GetString("archon.db.name")
	dsn := viper.GetString("archon.db.dsn")

	cli, err := orm.OpenEntClient(name, dsn)
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

func dbx(c *ent.Client) *orm.DBx {
	return &orm.DBx{Client: c}
}

func rootLog() *slog.Logger {
	rlogF := viper.GetString("archon.log.root.file")
	rlogL := viper.GetString("archon.log.root.level")
	logL := viper.GetString("archon.log.level")

	level := clog.Level(cmp.Or(rlogL, logL))
	log := clog.NewLog(rlogF, level)
	slog.SetDefault(log)
	return log
}

func serverLog() *slog.Logger {
	sLogF := viper.GetString("archon.log.server.file")
	sLogL := viper.GetString("archon.log.server.level")
	logL := viper.GetString("archon.log.level")

	level := clog.Level(cmp.Or(sLogL, logL))
	l := clog.NewLog(sLogF, level)
	return l
}
