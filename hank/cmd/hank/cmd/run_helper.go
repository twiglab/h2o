package cmd

import (
	"log"
	"log/slog"

	"github.com/spf13/viper"
	"github.com/twiglab/h2o/hank"
)

func buildDataLog() *slog.Logger {
	return hank.NewDataLog("logs/datalog.log")
}

func buildSender() hank.Sender {
	borker := viper.GetString("mqtt.borker")
	if borker == "" {
		log.Println("using logAction")
		return hank.NewLogAction()
	}
	cli, err := hank.NewMQTTClient("hank-plugin", borker)
	if err != nil {
		log.Fatal(err)
	}

	return hank.NewMQTTAction(cli)
}
