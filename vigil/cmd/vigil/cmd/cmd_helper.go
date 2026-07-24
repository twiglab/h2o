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
	"github.com/twiglab/h2o/pkg/common"
	"github.com/twiglab/h2o/vigil"
	"github.com/twiglab/h2o/vigil/orm"
	"github.com/twiglab/h2o/vigil/orm/ent"
	"github.com/twiglab/h2o/vigil/tsdb"
)

func mqttcli() mqtt.Client {
	broker := viper.GetString("vigil.mqtt.broker")
	if broker == "" {
		log.Fatalf("no broker")
	}
	cli, err := vigil.NewMQTTClient(vigil.CLIENT_ID, broker)
	if err != nil {
		log.Fatal(fmt.Errorf("mqttcli err: %w", err))
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
		log.Fatal(fmt.Errorf("ent err: %w", err))
	}
	return cli
}

func dbx(c *ent.Client) *orm.DBx {
	return &orm.DBx{Client: c}
}

func tdb() *tsdb.Schemaless {
	dsn := viper.GetString("vigil.tsdb.dsn")
	sch, err := tsdb.NewSchLe(dsn)
	if err != nil {
		log.Fatal(fmt.Errorf("SchLe err: %w", err))
	}
	return sch
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
