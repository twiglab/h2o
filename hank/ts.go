package hank

import (
	"bufio"
	"encoding/json/v2"
	"fmt"
	"log"
	"net"
	"time"
)

func (s *Server) RunX() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go xserve(conn)
	}
}

func xserve(conn net.Conn) {
	bio := bufio.NewReader(conn)
	for {
		bs, err := bio.ReadBytes('\n')
		if err != nil {
			log.Print(err)
			break
		}

		fmt.Println(string(bs))
		fmt.Println()
		fmt.Println(time.Now())
		fmt.Println()
		json.MarshalWrite(conn, OK)
	}
	fmt.Println("Break")
}
