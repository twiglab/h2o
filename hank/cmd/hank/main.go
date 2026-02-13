package main

import (
	"log"

	"net/http"
	_ "net/http/pprof"

	"github.com/twiglab/h2o/hank"
)

func main() {
	s := &hank.Server{
		Addr: "0.0.0.0:10004",
		Hub:  &hank.Hub{},
	}

	go http.ListenAndServe(":10007", nil)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
