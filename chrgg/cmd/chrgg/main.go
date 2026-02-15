package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/twiglab/h2o/chrgg"
)

var (
	broker string
	addr   string
	cid    string
)

func init() {
	flag.StringVar(&broker, "broker", "localhost:1883", "broker")
	flag.StringVar(&addr, "addr", ":10001", "addr")
	flag.StringVar(&cid, "cid", "chrgg", "cid")
}

func main() {
	flag.Parse()

	c, err := chrgg.NewMQTTClient(cid, broker)
	if err != nil {
		log.Fatal(err)
	}

	c.Subscribe("h2o/+/electricity", 0, chrgg.H)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
