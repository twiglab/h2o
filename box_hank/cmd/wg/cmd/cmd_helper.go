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
)

func rootLog() *slog.Logger {
	rlogF := viper.GetString("wg.log.root.file")
	rlogL := viper.GetString("wg.log.root.level")

	hlogL := viper.GetString("wg.log.level")
	level := clog.Level(cmp.Or(rlogL, hlogL))

	log := clog.NewLog(rlogF, level)
	slog.SetDefault(log)
	return log
}

func serverLog() *slog.Logger {
	sLogF := viper.GetString("wg.log.server.file")
	sLogL := viper.GetString("wg.log.server.level")
	hlogL := viper.GetString("wg.log.level")

	level := clog.Level(cmp.Or(sLogL, hlogL))
	l := clog.NewLog(sLogF, level)
	return l
}

func mqtt() *box.MQTTAction {
	broker := viper.GetString("wg.sender.mqtt.broker")
	clientID := box.ClientID("wg")
	log.Println("clientID", clientID)

	cli, err := box.NewMQTTClient(clientID, broker)
	if err != nil {
		log.Fatal(err)
	}
	return box.NewMQTTAction(cli)
}

func sender() box.Sender {
	use := viper.GetString("wg.sender.use")
	switch use {
	case "mqtt":
		log.Println("using mqtt")
		return mqtt()
	}
	log.Println("using logAction")
	return box.LogAction{}
}

func modbusClient() *modbus.ModbusClient {
	url := viper.GetString("wg.connect.modbus.rtu.url")
	log.Println("wg.connect.modbus.rtu.url:", url)

	dataBits := viper.GetUint("wg.connect.modbus.rtu.data")
	speed := viper.GetUint("wg.connect.modbus.rtu.speed")
	stopBits := viper.GetUint("wg.connect.modbus.rtu.stop")

	cfg := &modbus.ClientConfiguration{
		URL:      url,
		Speed:    speed,
		DataBits: dataBits,
		StopBits: stopBits,
		Timeout:  3 * time.Second,
	}
	cli, err := box.NewModbusClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}
