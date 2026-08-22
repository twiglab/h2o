package cmd

import (
	"cmp"
	"log"
	"log/slog"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/box"
	"github.com/twiglab/h2o/clog"
	"github.com/twiglab/h2o/clog/wal"
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
		log.Fatalln("wal file is null. ***MUST*** set ocg.wal.file")
	}
	log.Println("wal file:", logf)
	return wal.New(wal.Conf{Filename: logf})
}

func mqtt() *box.MQTTAction {
	broker := viper.GetString("ocg.sender.mqtt.broker")
	clientID := box.ClientID("ocg")
	log.Println("clientID", clientID)

	cli, err := box.NewMQTTClient(clientID, broker)
	if err != nil {
		log.Fatal(err)
	}
	return box.NewMQTTAction(cli)
}

func sender() box.Sender {
	use := viper.GetString("ocg.sender.use")
	switch use {
	case "mqtt":
		log.Println("using mqtt")
		return mqtt()
	}
	log.Println("using logAction")
	return box.LogAction{}
}

func modbusClient() *modbus.ModbusClient {
	url := viper.GetString("ocg.connect.modbus.tcp.url")
	log.Println("ocg.connect.modbus.tcp.url:", url)
	cfg := &modbus.ClientConfiguration{
		URL:     url,
		Timeout: 3 * time.Second,
	}
	cli, err := box.NewModbusClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}
