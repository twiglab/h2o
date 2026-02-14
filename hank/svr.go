package hank

import (
	"bufio"
	"context"
	"log/slog"
	"net"

	"github.com/cloudwego/netpoll"
	"github.com/google/uuid"
)

type connKey struct {
	name string
}

var ck = connKey{"__conn_id"}

type cid struct {
	s  *Server
	id string
}

func (c cid) String() string {
	return c.id
}

type Server struct {
	Addr string
	Hub  *Hub
	Enh  *Enh

	BaseCtx func(net.Conn) context.Context

	Logger *slog.Logger
}

func (s *Server) RunAt(l net.Listener) error {
	loop, err := netpoll.NewEventLoop(
		at(serve),

		netpoll.WithOnDisconnect(func(ctx context.Context, conn netpoll.Connection) {
			id := fromCtx[*cid](ctx, ck)
			s.Logger.DebugContext(ctx, "onConnect",
				slog.String("remoteAddr", conn.RemoteAddr().String()),
				slog.String("cid", id.String()),
			)
		}),

		netpoll.WithOnConnect(func(ctx context.Context, conn netpoll.Connection) context.Context {
			sk := &cid{s: s, id: uuid.NewString()}

			s.Logger.DebugContext(ctx, "onConnect",
				slog.String("remoteAddr", conn.RemoteAddr().String()),
				slog.String("cid", sk.String()),
			)

			return context.WithValue(ctx, ck, sk)
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

func at(f func(context.Context, net.Conn, *Server) error) netpoll.OnRequest {
	return func(ctx context.Context, conn netpoll.Connection) error {
		id := fromCtx[*cid](ctx, ck)
		return f(ctx, conn, id.s)
	}
}

func fromCtx[T any](ctx context.Context, key any) T {
	return ctx.Value(key).(T)
}

func serve(ctx context.Context, conn net.Conn, s *Server) error {
	sc := bufio.NewScanner(conn)
	sk := fromCtx[*cid](ctx, ck)
	for sc.Scan() {
		var sd SyncData
		if err := unmarshal(sc.Bytes(), &sd); err != nil {
			s.Logger.ErrorContext(ctx, "unmarshal SyncData error",
				slog.String("remoteAddr", conn.RemoteAddr().String()),
				slog.String("cid", sk.String()),
				slog.Any("error", err),
			)
			continue
		}

		if sd.Type == TypeRate {
			err := marshalWrite(conn, ErrNoRate)
			s.Logger.InfoContext(ctx, "rate type", slog.String("type", sd.Type), slog.Any("error", err))
			continue
		}

		switch sd.Type {
		case TypeDeviceData:
			go doDeviceData(ctx, sd, s)
		case TypeDeviceStatus:
			go doDeviceStatus(ctx, sd, s)
		default:
			s.Logger.InfoContext(ctx, "ignore type", slog.String("type", sd.Type))
		}

		if err := marshalWrite(conn, OK); err != nil {
			s.Logger.ErrorContext(ctx, "unmarshalWriter OK error", slog.Any("error", err))
		}
	}
	return sc.Err()
}

func doDeviceData(ctx context.Context, sd SyncData, s *Server) {
	var ddl DeviceDataList
	if err := unmarshal(sd.Data, &ddl); err != nil {
		s.Logger.ErrorContext(ctx, "unmarshal deviceDataList error", slog.Any("error", err))
		return
	}

	for _, dd := range ddl {
		if err := s.Hub.HandleDeviceData(ctx, s.Enh.Convert(dd)); err != nil {
			s.Logger.ErrorContext(ctx, "handleDeviceData error", slog.Any("error", err))
		}
	}
}

func doDeviceStatus(ctx context.Context, sd SyncData, s *Server) {
	var dsl DeviceStatusList
	if err := unmarshal(sd.Data, &dsl); err != nil {
		s.Logger.ErrorContext(ctx, "unmarshal deviceStatusList error", slog.Any("error", err))
		return
	}
	for _, ds := range dsl {
		if err := s.Hub.HandleDeviceStatus(ctx, ds); err != nil {
			s.Logger.ErrorContext(ctx, "handleDeviceStatus error", slog.Any("error", err))
		}
	}
}
