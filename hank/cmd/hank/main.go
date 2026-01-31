package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"

	hank "github.com/twiglab/h2o/octopus"
)

var (
	broker  string
	logPath string
)

func init() {
	flag.StringVar(&broker, "broker", "127.0.0.1:1088", "broker")
	flag.StringVar(&logPath, "logpath", ".", "logpath")
}

func main() {
	flag.Parse()

	fmt.Println(broker)
	fmt.Println(logPath)

	datalog := filepath.Join(logPath, "data.log")
	infolog := filepath.Join(logPath, "info.log")

	dl := hank.NewLog(datalog, slog.LevelDebug)
	il := hank.NewLog(infolog, slog.LevelDebug)

	hub := &hank.Hub{
		DataLog: dl,
		InfoLog: il,
		Enh:     &hank.Enh{},
		Sender:  mqtt(broker),
	}
	h := hank.Handle(hub)

	http.HandleFunc("/sync", h)

	if err := http.ListenAndServe("0.0.0.0:20001", nil); err != nil {
		log.Fatal(err)
	}
}

func mqtt(borker string) *hank.MQTTAction {
	c, err := hank.NewMQTTClient("hank", borker)
	if err != nil {
		log.Fatal(err)
	}

	return hank.NewMQTTAction(c)
}
