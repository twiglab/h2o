package hank

import (
	"bufio"
	"context"
	"io"
	"log"
	"log/slog"
	"net"

	"github.com/cloudwego/netpoll"
	"github.com/google/uuid"
)

type connKey struct {
	name string
}

var sidKey = connKey{"__conn_sid_key__"}

type sid struct {
	s  *Server
	id string
}

func (s sid) String() string {
	return s.id
}

type Server struct {
	Addr string
	Hub  *Hub
	Enh  *Enh

	BaseCtx func(netpoll.Connection) context.Context

	Logger *slog.Logger
}

func (s *Server) RunAt(l net.Listener) error {
	loop, err := netpoll.NewEventLoop(
		at(serve),

		netpoll.WithOnDisconnect(func(ctx context.Context, conn netpoll.Connection) {
			log.Println("disconnect ... ", conn.RemoteAddr(), "sid ...", ctx.Value(sidKey))
		}),

		netpoll.WithOnConnect(func(ctx context.Context, conn netpoll.Connection) context.Context {
			_ = conn.AddCloseCallback(func(conn netpoll.Connection) error {
				log.Println("closing ... ", conn.RemoteAddr())
				return nil
			})
			sk := &sid{s: s, id: uuid.NewString()}
			log.Println("connect ... ", conn.RemoteAddr(), "sid...", sk)
			return context.WithValue(ctx, sidKey, sk)
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

func at(f func(context.Context, io.ReadWriteCloser, *Server) error) netpoll.OnRequest {
	return func(ctx context.Context, conn netpoll.Connection) error {
		v := ctx.Value(sidKey).(*sid)
		return f(ctx, conn, v.s)
	}
}

func serve(ctx context.Context, conn io.ReadWriteCloser, s *Server) error {
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		var sd SyncData
		if err := unmarshal(sc.Bytes(), &sd); err != nil {
			log.Println(err, "SD")
			continue
		}

		if sd.Type == TypeRate {
			_ = marshalWrite(conn, ErrNoRate)
			continue
		}

		switch sd.Type {
		case TypeDeviceData:
			go doDeviceData(ctx, sd, s)
		case TypeDeviceStatus:
			go doDeviceStatus(ctx, sd, s)
		default:
			log.Println("ignore type = ", sd.Type)
		}

		if err := marshalWrite(conn, OK); err != nil {
			log.Println(err, "marshalWriter")
		}
	}
	return sc.Err()
}

func doDeviceData(ctx context.Context, sd SyncData, s *Server) {
	var ddl DeviceDataList
	if err := unmarshal(sd.Data, &ddl); err != nil {
		log.Println(err, "ddl")
		return
	}

	for _, dd := range ddl {
		if err := s.Hub.HandleDeviceData(ctx, s.Enh.Convert(dd)); err != nil {
			log.Println(err)
		}
	}
}

func doDeviceStatus(ctx context.Context, sd SyncData, s *Server) {
	var dsl DeviceStatusList
	if err := unmarshal(sd.Data, &dsl); err != nil {
		log.Println(err, "dsl")
		return
	}
	for _, ds := range dsl {
		if err := s.Hub.HandleDeviceStatus(ctx, ds); err != nil {
			log.Println(err)
		}
	}
}
