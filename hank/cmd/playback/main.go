package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	filename string
	addr     string
)

func init() {
	flag.StringVar(&filename, "file", "", "file")
	flag.StringVar(&addr, "addr", "127.0.0.1:10004", "addr")
}

func main() {

	flag.Parse()

	if filename == "" {
		log.Fatalln("filename is null")
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln("连接服务器失败:", err)
	}
	defer conn.Close()
	fmt.Println("已连接到服务器")

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		w.Write(sc.Bytes())
		w.WriteByte('\n')
		w.Flush()
		r.ReadBytes('\n')
	}
}
