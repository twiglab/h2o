package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	hank "github.com/twiglab/h2o/octopus"
)

var (
	borker     string
	logPath    string
	serverPath string
)

func init() {
	flag.StringVar(&borker, "borker", "127.0.0.1:1088", "borker")
	flag.StringVar(&logPath, "logpath", ".", "logpath")
	flag.StringVar(&serverPath, "ctxpath", "/sync", "ctxpath")
}

func main() {
	flag.Parse()

	fmt.Println(borker)
	fmt.Println(logPath)
	fmt.Println(serverPath)

	dl := hank.NewLog("data.log", slog.LevelDebug)
	il := hank.NewLog("info.log", slog.LevelDebug)

	hub := &hank.Hub{
		DataLog: dl,
		InfoLog: il,
	}
	h := hank.Handle(hub)

	http.HandleFunc(serverPath, h)

	if err := http.ListenAndServe("0.0.0.0:20000", nil); err != nil {
		log.Fatal(err)
	}
}
