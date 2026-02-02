package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/twiglab/h2o/hank"
)

var (
	broker  string
	logPath string
)

func init() {
	flag.StringVar(&broker, "broker", "localhost:1883", "broker")
	flag.StringVar(&logPath, "logpath", "logs", "logpath")
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
	if borker == "" || borker == "nil" || borker == "x" {
		return nil
	}

	c, err := hank.NewMQTTClient("hank-plugin", borker)
	if err != nil {
		log.Fatal(err)
	}

	return hank.NewMQTTAction(c)
}
