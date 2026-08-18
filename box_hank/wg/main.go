package main

import (
	"fmt"
	"log"
	"time"

	"github.com/simonvetter/modbus"
)

func main() {
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:      "rtu:///dev/tty0",
		Speed:    9600,               // default
		DataBits: 8,                  // default, optional
		Parity:   modbus.PARITY_NONE, // default, optional
		StopBits: 1,                  // default if no parity, optional
		Timeout:  1 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := client.Open(); err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	client.SetUnitId(21)

	val, err := client.ReadUint32(0, modbus.INPUT_REGISTER)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(val)
}
