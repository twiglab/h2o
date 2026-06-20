package hank

import (
	"bufio"
	"context"
	"log/slog"
	"math/rand/v2"
	"net"

	"github.com/cloudwego/netpoll"
)

type connKey struct {
	name string
}

var ck = connKey{"__conn_id"}

type cid struct {
	s  *Server
	id int64
}

type Server struct {
	Addr string
	Hub  *Hub
	Enh  *Enh

	BaseCtx func(net.Conn) context.Context

	Logger *slog.Logger

	PlayBack *PlayBack
}

func (s *Server) RunAt(l net.Listener) error {
	loop, err := netpoll.NewEventLoop(
		at(serve),

		netpoll.WithOnConnect(func(ctx context.Context, conn netpoll.Connection) context.Context {
			sk := &cid{s: s, id: rand.Int64()}
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

	slog.DebugContext(ctx, "serve",
		slog.String("remoteAddr", conn.RemoteAddr().String()),
		slog.Int64("cid", sk.id),
	)

	for sc.Scan() {
		s.PlayBack.Record(ctx, sc.Text())

		var sd SyncData
		if err := unmarshal(sc.Bytes(), &sd); err != nil {
			s.Logger.ErrorContext(ctx, "unmarshal SyncData error",
				slog.String("remoteAddr", conn.RemoteAddr().String()),
				slog.Int64("cid", sk.id),
				slog.Any("error", err),
			)
			continue
		}

		if err := writeReturn(conn, OK); err != nil {
			// 和对方确认，网关发送完毕数据2s后断开，但是经过实际测试，网关并没有2s的延时，应该是发送完毕就直接断开了
			// 另外对方答复不返回ok会导致后续再次发送，实际运行也没有发现再次发送的情况
			// 另外实际业务处理，也不可能达到2s之久
			// 综上所述，保留这个日志，并将等级由于Error改为Debug
			// 后续待查
			// TODO
			s.Logger.DebugContext(ctx, "unmarshalRetrun OK error", slog.Any("error", err))
		}

		if sd.Type == TypeRate {
			err := writeReturn(conn, ErrNoRate)
			s.Logger.InfoContext(ctx, "rate type", slog.String("type", sd.Type), slog.Any("error", err))
			continue
		}

		switch sd.Type {
		case TypeDeviceData:
			doDeviceData(ctx, sd, s)
		case TypeDeviceStatus:
			doDeviceStatus(ctx, sd, s)
		default:
			s.Logger.InfoContext(ctx, "ignore type", slog.String("type", sd.Type))
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
		switch dd.Type {
		case ELECTRICITY:
			em, err := s.Enh.ToElecty(dd)
			if err != nil {
				s.Logger.ErrorContext(ctx, "Enh.ToElecty error", slog.Any("raw", dd), slog.Any("error", err))
				return
			}
			if err := s.Hub.HandleElectricity(ctx, em); err != nil {
				s.Logger.ErrorContext(ctx, "handleElectricity error", slog.Any("data", em), slog.Any("error", err))
				return
			}
		case WATER:
			wm, err := s.Enh.ToWater(dd)
			if err != nil {
				s.Logger.ErrorContext(ctx, "Enh.ToWater error", slog.Any("raw", dd), slog.Any("error", err))
				return
			}
			if err := s.Hub.HandleWater(ctx, wm); err != nil {
				s.Logger.ErrorContext(ctx, "handleWater error", slog.Any("data", wm), slog.Any("error", err))
				return
			}
		default:
			s.Logger.ErrorContext(ctx, "unknow device type", slog.String("type", dd.Type))
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
