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
	"github.com/google/uuid"
)

type connKey struct {
	name string
}

var ckey = connKey{"__conn_key__"}

type Server struct {
	Addr string
	Hub  *Hub
	Enh  *Enh

	BaseCtx func(netpoll.Connection) context.Context
}

func (s *Server) RunAt(l net.Listener) error {
	loop, err := netpoll.NewEventLoop(
		at(s.serv),

		netpoll.WithOnDisconnect(func(ctx context.Context, conn netpoll.Connection) {
			fmt.Println("disconnect ..... ", conn.RemoteAddr(), "key ...", ctx.Value(ckey))
		}),

		netpoll.WithOnConnect(func(ctx context.Context, conn netpoll.Connection) context.Context {
			id := uuid.NewString()
			fmt.Println("connect... ", conn.RemoteAddr(), "key...", id)
			fmt.Println("enh...", s.Enh)
			return context.WithValue(ctx, ckey, id)
		}),

		netpoll.WithOnPrepare(func(conn netpoll.Connection) context.Context {
			if s.BaseCtx != nil {
				return s.BaseCtx(conn)
			}
			return context.Background()
		}),
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

func (s *Server) serv(ctx context.Context, rwc io.ReadWriteCloser) error {
	sc := bufio.NewScanner(rwc)
	for sc.Scan() {
		d := &SyncData{}
		if err := json.Unmarshal(sc.Bytes(), d); err != nil {
			log.Println(err, "SD")
			return err
		}

		if d.Type == TypeRate {
			_ = json.MarshalWrite(rwc, ErrNoRate)
			continue
		}

		if d.Type == TypeDeviceData {
			var ddl DeviceDataList
			if err := json.Unmarshal(d.Data, &ddl); err != nil {
				log.Println(err, "ddl")
				continue
			}
			go s.doDeviceDataList(ctx, ddl)
		}

		if d.Type == TypeDeviceStatus {
			var dsl DeviceStatusList
			if err := json.Unmarshal(d.Data, &dsl); err != nil {
				log.Println(err, "dsl")
				continue
			}
			go s.doDeviceStatusList(ctx, dsl)
		}

		_ = json.MarshalWrite(rwc, OK)
	}
	return sc.Err()
}

func (s *Server) doDeviceDataList(ctx context.Context, ddl DeviceDataList) {
	for _, dd := range ddl {
		if err := s.Hub.HandleDeviceData(ctx, s.Enh.Convert(dd)); err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) doDeviceStatusList(ctx context.Context, dsl DeviceStatusList) {
	for _, ds := range dsl {
		if err := s.Hub.HandleDeviceStatus(ctx, ds); err != nil {
			log.Println(err)
		}
	}
}

func at(f func(context.Context, io.ReadWriteCloser) error) netpoll.OnRequest {
	return func(ctx context.Context, conn netpoll.Connection) error {
		return f(ctx, conn)
	}
}
