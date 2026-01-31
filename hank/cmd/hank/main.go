package main

import (
	"log"
	"log/slog"
	"net/http"

	hank "github.com/twiglab/h2o/octopus"
)

func main() {

	dl := hank.NewLog("data.log", slog.LevelDebug)
	il := hank.NewLog("info.log", slog.LevelDebug)

	hub := &hank.Hub{
		DataLog: dl,
		InfoLog: il,
	}
	h := hank.Handle(hub)

	http.HandleFunc("/sync", h)

	if err := http.ListenAndServe("0.0.0.0:20000", nil); err != nil {
		log.Fatal(err)
	}
}
