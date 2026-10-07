package cmd

import (
	"cmp"
	"fmt"
	"log"
	"log/slog"

	"github.com/go-sphere/wecom-aibot-go-sdk/aibot"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/clog/wal"
	"github.com/twiglab/h2o/dbcli"
	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/topup/wecom"
)

func webaddr() string {
	addr := viper.GetString("topup.web.addr")
	return cmp.Or(addr, ":10009")
}

func entcli() *ent.Client {
	name := viper.GetString("topup.db.name")
	dsn := viper.GetString("topup.db.dsn")

	cli, err := dbcli.OpenEntClient(name, dsn)
	if err != nil {
		log.Fatal(fmt.Errorf("ent err: %w", err))
	}
	return cli
}

func wsclinet() *aibot.WSClient {
	botID := viper.GetString("topup.wecom.bot.bot_id")
	botSecret := viper.GetString("topup.wecom.bot.bot_secret")
	if botID == "" || botSecret == "" {
		log.Fatal("请设置 BOT_ID 和 BOT_SECRET")
	}

	return wecom.NewWsClient(botID, botSecret)
}

func fixGroup() *wecom.FixUserGroup {
	chatID := viper.GetString("topup.auth.fix.chat_id")
	users := viper.GetStringSlice("topup.auth.fix.users")
	return &wecom.FixUserGroup{
		Users:  users,
		ChatID: chatID,
	}
}

func rootLog() *slog.Logger {
	rlogF := viper.GetString("vigil.log.root.file")
	rlogL := viper.GetString("vigil.log.root.level")
	logL := viper.GetString("vigil.log.level")

	level := clog.Level(cmp.Or(rlogL, logL))
	log := clog.NewLog(rlogF, level)
	slog.SetDefault(log)
	return log
}

func serverLog() *slog.Logger {
	sLogF := viper.GetString("vigil.log.server.file")
	sLogL := viper.GetString("vigil.log.server.level")
	logL := viper.GetString("vigil.log.level")

	level := clog.Level(cmp.Or(sLogL, logL))
	l := clog.NewLog(sLogF, level)
	return l
}

func wallog() *wal.WAL {
	logf := viper.GetString("vigil.wal.file")
	if logf == "" {
		log.Fatalln("wal file is null. ***MUST*** set vigil.wal.file")
	}
	log.Println("wal file:", logf)
	return wal.New(wal.Conf{Filename: logf})
}
