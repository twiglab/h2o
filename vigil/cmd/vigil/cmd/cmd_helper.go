package cmd

import (
	"cmp"
	"log"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	vlog "github.com/twiglab/h2o/log"
	"github.com/twiglab/h2o/pkg/common"
	"github.com/twiglab/h2o/vigil"
	"github.com/twiglab/h2o/vigil/orm"
	"github.com/twiglab/h2o/vigil/orm/ent"
)

func logLevel(s string) slog.Level {
	switch s {
	case "debug", "DEBUG":
		return slog.LevelDebug
	case "info", "INFO":
		return slog.LevelInfo
	case "error", "ERROR":
		return slog.LevelError
	case "warn", "WARN":
		return slog.LevelWarn
	}
	return slog.LevelInfo
}

func mqttcli() mqtt.Client {
	broker := viper.GetString("vigil.mqtt.broker")
	if broker == "" {
		log.Fatalf("no broker")
	}
	cli, err := vigil.NewMQTTClient(vigil.CLIENT_ID, broker)
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

func webaddr() string {
	addr := viper.GetString("vigil.web.addr")
	return cmp.Or(addr, ":10003")
}

func topics() map[string]byte {
	return map[string]byte{
		common.WaterTopic:       0x01,
		common.ElectricityTopic: 0x01,
		common.GasTopic:         0x01,
	}
}

func entcli() *ent.Client {
	name := viper.GetString("vigil.db.name")
	dsn := viper.GetString("vigil.db.dsn")

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
	rlogF := viper.GetString("log.root.file")
	rlogL := viper.GetString("log.root.level")
	logL := viper.GetString("log.level")

	level := logLevel(cmp.Or(rlogL, logL))
	log := vlog.NewLog(rlogF, level)
	slog.SetDefault(log)
	return log
}

func serverLog() *slog.Logger {
	sLogF := viper.GetString("log.server.file")
	sLogL := viper.GetString("log.server.level")
	logL := viper.GetString("log.level")

	level := logLevel(cmp.Or(sLogL, logL))
	l := vlog.NewLog(sLogF, level)
	return l
}
