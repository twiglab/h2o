package hank

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"encoding/json/v2"

	"github.com/cloudwego/netpoll"
)

type Server struct {
	Addr string
	Hub  *Hub
	Enh  *Enh

	BaseCtx func(netpoll.Connection) context.Context
	ConnCtx func(context.Context, netpoll.Connection) context.Context
}

func (s *Server) RunAt(l net.Listener) error {
	loop, err := netpoll.NewEventLoop(
		at(s.serv),

		netpoll.WithOnDisconnect(func(ctx context.Context, conn netpoll.Connection) {
			fmt.Println("Dis .....")
		}),

		netpoll.WithOnConnect(func(ctx context.Context, conn netpoll.Connection) context.Context {
			if s.ConnCtx != nil {
				return s.ConnCtx(ctx, conn)
			}
			return ctx
		}),

		netpoll.WithOnPrepare(func(conn netpoll.Connection) context.Context {
			if s.BaseCtx != nil {
				return s.BaseCtx(conn)
			}
			return context.Background()
		}),

		//netpoll.WithReadTimeout(30*time.Second),
	)
	if err != nil {
		return err
	}
	return loop.Serve(l)
}

func (s *Server) Run() error {
	ln, err := netpoll.CreateListener("tcp4", s.Addr)
	if err != nil {
		return err
	}
	return s.RunAt(ln)
}

func (s *Server) serv(ctx context.Context, rw io.ReadWriteCloser) error {
	sc := bufio.NewScanner(rw)
	for sc.Scan() {
		d := &SyncData{}
		if err := json.Unmarshal(sc.Bytes(), d); err != nil {
			log.Println(err, "SD")
			return err
		}

		if d.Type == TypeRate {
			_ = json.MarshalWrite(rw, ErrNoRate)
			continue
		}

		if d.Type == TypeDeviceData {
			var ddl DeviceDataList
			if err := json.Unmarshal(d.Data, &ddl); err != nil {
				log.Println(err, "ddl")
				continue
			}

			doHandleDeviceDataList(ctx, ddl, s.Hub)
		}

		if d.Type == TypeDeviceStatus {
			var dsl DeviceStatusList
			if err := json.Unmarshal(d.Data, &dsl); err != nil {
				log.Println(err, "dsl")
				continue
			}
			doHandleDeviceStatusList(ctx, dsl, s.Hub)
		}

		_ = json.MarshalWrite(rw, OK)
	}
	return sc.Err()
}

func at(f func(context.Context, io.ReadWriteCloser) error) netpoll.OnRequest {
	return func(ctx context.Context, conn netpoll.Connection) error {
		return f(ctx, conn)
	}
}
