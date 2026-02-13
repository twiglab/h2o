package main

import (
	"log"

	"github.com/twiglab/octopus/hank"
)

func main() {
	s := &hank.Server{
		Addr: "0.0.0.0:10004",
		Hub:  &hank.Hub{},
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
