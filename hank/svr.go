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

type svr struct {
	s  *Server
	id string
}

func (s svr) String() string {
	return s.id
}

type Server struct {
	Addr string
	Hub  *Hub
	Enh  *Enh

	BaseCtx func(netpoll.Connection) context.Context
}

func (s *Server) RunAt(l net.Listener) error {
	loop, err := netpoll.NewEventLoop(
		at(serve),

		netpoll.WithOnDisconnect(func(ctx context.Context, conn netpoll.Connection) {
			fmt.Println("disconnect ..... ", conn.RemoteAddr(), "key ...", ctx.Value(ckey))
		}),

		netpoll.WithOnConnect(func(ctx context.Context, conn netpoll.Connection) context.Context {
			_ = conn.AddCloseCallback(func(conn netpoll.Connection) error {
				fmt.Println("closing ...")
				return nil
			})
			id := uuid.NewString()
			fmt.Println("connect... ", conn.RemoteAddr(), "key...", id)
			sk := &svr{s: s, id: id}
			return context.WithValue(ctx, ckey, sk)
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
	ln, err := netpoll.CreateListener("tcp", s.Addr)
	if err != nil {
		return err
	}
	return s.RunAt(ln)
}

func doDataList(ctx context.Context, ddl DeviceDataList, s *Server) {
	for _, dd := range ddl {
		if err := s.Hub.HandleDeviceData(ctx, s.Enh.Convert(dd)); err != nil {
			log.Println(err)
		}
	}
}

func doStatusList(ctx context.Context, dsl DeviceStatusList, s *Server) {
	for _, ds := range dsl {
		if err := s.Hub.HandleDeviceStatus(ctx, ds); err != nil {
			log.Println(err)
		}
	}
}

func at(f func(context.Context, io.ReadWriteCloser, *Server) error) netpoll.OnRequest {
	return func(ctx context.Context, conn netpoll.Connection) error {
		v := ctx.Value(ckey).(*svr)
		return f(ctx, conn, v.s)
	}
}

func serve(ctx context.Context, conn io.ReadWriteCloser, s *Server) error {
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		var d SyncData
		if err := json.Unmarshal(sc.Bytes(), &d); err != nil {
			log.Println(err, "SD")
			continue
		}

		if d.Type == TypeRate {
			log.Println(d.Type)
			_ = json.MarshalWrite(conn, ErrNoRate)
			continue
		}

		if d.Type == TypeDeviceData {
			var ddl DeviceDataList
			if err := json.Unmarshal(d.Data, &ddl); err != nil {
				log.Println(err, "ddl")
				continue
			}
			go doDataList(ctx, ddl, s)
		}

		if d.Type == TypeDeviceStatus {
			var dsl DeviceStatusList
			if err := json.Unmarshal(d.Data, &dsl); err != nil {
				log.Println(err, "dsl")
				continue
			}
			go doStatusList(ctx, dsl, s)
		}

		if err := json.MarshalWrite(conn, OK); err != nil {
			log.Print(err, "marshalWriter")
		}
	}
	return sc.Err()
}
