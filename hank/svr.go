package hank

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"encoding/json/jsontext"
	"encoding/json/v2"

	"github.com/cloudwego/netpoll"
)

func x(ctx context.Context, conn netpoll.Connection) error {
	dec := jsontext.NewDecoder(conn)
	for {
		var d SyncData
		if err := json.UnmarshalDecode(dec, &d); err != nil {
			return err
		}
		fmt.Println(d.Type)
	}
}

func main() {
	l, err := netpoll.CreateListener("tcp4", "0.0.0.0:10004")
	if err != nil {
		log.Fatal(err)
	}

	loop, err := netpoll.NewEventLoop(x)
	if err != nil {
		log.Fatal(err)
	}

	if err := loop.Serve(l); err != nil {
		log.Fatal(err)
	}
}

type Job func(context.Context, io.ReadWriteCloser) error

type Server struct {
	Addr string
	Job  Job
	loop netpoll.EventLoop
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp4", s.Addr)
	if err != nil {
		return err
	}

	l, err := netpoll.NewEventLoop(nil)
	if err != nil {
		return err
	}
	s.loop = l

	if err := s.loop.Serve(ln); err != nil {
		return err
	}
}
