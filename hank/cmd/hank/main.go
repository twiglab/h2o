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
	addr    string
)

func init() {
	// flag.StringVar(&broker, "broker", "localhost:1883", "broker")
	flag.StringVar(&broker, "broker", "log", "broker")
	flag.StringVar(&logPath, "logpath", "logs", "logpath")
	flag.StringVar(&addr, "addr", "0.0.0.0:10007", "addr")
}

func main() {
	flag.Parse()

	fmt.Println(broker)
	fmt.Println(logPath)
	fmt.Println(addr)

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

	http.HandleFunc(hank.SERVER_URL_PATH, h)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func mqtt(borker string) hank.Sender {
	if borker == "" || borker == "log" || borker == "x" {
		return hank.NewLogAction()
	}

	c, err := hank.NewMQTTClient("hank-plugin", borker)
	if err != nil {
		log.Fatal(err)
	}

	return hank.NewMQTTAction(c)
}
