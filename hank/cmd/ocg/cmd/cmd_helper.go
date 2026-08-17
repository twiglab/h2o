package cmd

import (
	"cmp"
	"log"
	"log/slog"

	"github.com/simonvetter/modbus"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/clog/wal"
	"github.com/twiglab/h2o/hank"
	"github.com/twiglab/h2o/hank/cmd/ocg/oc"
)

func rootLog() *slog.Logger {
	rlogF := viper.GetString("ocg.log.root.file")
	rlogL := viper.GetString("ocg.log.root.level")

	hlogL := viper.GetString("ocg.log.level")
	level := clog.Level(cmp.Or(rlogL, hlogL))

	log := clog.NewLog(rlogF, level)
	slog.SetDefault(log)
	return log
}

func serverLog() *slog.Logger {
	sLogF := viper.GetString("ocg.log.server.file")
	sLogL := viper.GetString("ocg.log.server.level")
	hlogL := viper.GetString("ocg.log.level")

	level := clog.Level(cmp.Or(sLogL, hlogL))
	l := clog.NewLog(sLogF, level)
	return l
}

func wallog() *wal.WAL {
	logf := viper.GetString("ocg.wal.file")
	if logf == "" {
		log.Fatalln("wal file is null. ***MUST*** set hank.wal.file")
	}
	log.Println("wal file:", logf)
	return wal.New(wal.Conf{Filename: logf})
}

func mqtt() *hank.MQTTAction {
	broker := viper.GetString("ocg.sender.mqtt.broker")
	clientID := oc.ClientID("ocg")
	log.Println("clientID", clientID)

	cli, err := hank.NewMQTTClient(clientID, broker)
	if err != nil {
		log.Fatal(err)
	}
	return hank.NewMQTTAction(cli)
}

func sender() hank.Sender {
	use := viper.GetString("ocg.sender.use")
	switch use {
	case "mqtt":
		log.Println("using mqtt")
		return mqtt()
	}
	log.Println("using logAction")
	return hank.LogAction{}
}

func modbusClient(url string) *modbus.ModbusClient {
	cli, err := oc.NewModbusClient(url)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}
