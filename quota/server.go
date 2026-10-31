package quota

import (
	"cmp"
	"context"
	"log"
	"log/slog"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/clog/wal"
	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/proto"
)

const GeneralDataTopic = "h2o/data/#"

type ChargeServer struct {
	DBx    *DBx
	Logger *slog.Logger
	WAL    *wal.WAL

	sender Sender
}

func (s *ChargeServer) Run() mqtt.OnConnectHandler {
	return func(client mqtt.Client) {
		s.sender = NewMQTTAction(client)

		tok := client.Subscribe(GeneralDataTopic, 0x01, s.msgHandle())
		if !tok.WaitTimeout(10*time.Second) || tok.Error() != nil {
			log.Printf("resubscribe failed: %v", tok.Error())
		}
	}
}

func (s *ChargeServer) msgHandle() mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			msg.Ack()
			return
		}

		defer msg.Ack()

		// 协程处理
		go func(m mqtt.Message) {
			var em Meter
			if err := em.UnmarshalBinary(m.Payload()); err != nil {
				s.Logger.Error("unmarshal error", slog.Any("error", err))
				return
			}
			if err := s.charge(context.Background(), em); err != nil {
				s.Logger.Error("charge error", slog.Any("raw", em), slog.Any("error", err))
			}
		}(msg)

	}
}

func (s *ChargeServer) doStatusOff(ctx context.Context, md Meter, vc *ent.Top) error {
	if cmp.Less(md.Data.DataValue, vc.Top) {
		// 在断开状态，小于限额，发送合闸消息，开
		s.WAL.WriteLogContext(ctx,
			wal.String("top.code", vc.Code),
			wal.Int64("top", vc.Top), wal.Int64("dataValue", md.Data.DataValue),
			wal.String("code", md.Code), wal.String("type", md.Type),
			wal.String("onoff", "on"), // 超出限额，正常合闸操作
			wal.String("msg", "正常合闸"),
		)

		_ = vc.Update().SetStatus(STATUS_BEGIN).Exec(ctx)

		ot := newOnOffMessage(md, vc, proto.ON)
		return s.sender.SendData(ctx, ot)
	}
	return nil
}

func (s *ChargeServer) doStatusOn(ctx context.Context, md Meter, vc *ent.Top) error {
	if cmp.Less(vc.Top, md.Data.DataValue) {
		// 超额, 拉闸断开
		if vc.Status == STATUS_BEGIN {
			// 限额记录正常，执行拉闸操作
			s.WAL.WriteLogContext(ctx,
				wal.String("top.code", vc.Code),
				wal.Int64("top", vc.Top), wal.Int64("dataValue", md.Data.DataValue),
				wal.String("code", md.Code), wal.String("type", md.Type),
				wal.String("onoff", "off"), // 超出限额，正常拉闸操作
				wal.String("msg", "超出限额"),
			)

			_ = vc.Update().
				SetStatus(STATUS_END).
				SetEndTime(time.Now()).         // 计费结束时间
				SetEndStock(md.Data.DataValue). // 计费结束的表显
				Exec(ctx)                       // 计费结束

			// 断开
			ot := newOnOffMessage(md, vc, proto.OFF)
			return s.sender.SendData(ctx, ot)
		} else {
			// 非正常状态，疑似数据有非法修改
			s.Logger.WarnContext(ctx, "illegal status",
				slog.String("top.code", vc.Code),
				slog.Int64("top", vc.Top),
				slog.Int64("dataValue", md.Data.DataValue),
				slog.Int("status", vc.Status),
				slog.Any("meter", md), slog.Any("topRec", vc))

			// 把状态改成 begin 等待下一轮判断
			_ = vc.Update().SetStatus(STATUS_BEGIN).Exec(ctx)
		}
	}

	return nil
}

func (s *ChargeServer) charge(ctx context.Context, md Meter) error {
	if md.Data.OptStatus == proto.OPT_STATUS_UNKNOW {
		s.Logger.DebugContext(ctx, "status unknow", slog.Any("meter", md))
		return nil
	}

	vc, err := s.DBx.LoadLast(ctx, md.Code, md.Type)
	if ent.IsNotFound(err) {
		// 没找到限额，默认不计费
		s.Logger.DebugContext(ctx, "record not found", slog.Any("meter", md))
		return nil
	}

	// 其他错误
	if ent.MaskNotFound(err) != nil {
		s.Logger.ErrorContext(ctx, "loadLast error", slog.Any("meter", md), slog.Any("error", err))
		return err
	}

	// 注意： 即便状态为 STATUS_END 也是要继续执行的
	// 是否开关，完全由限额决定
	switch md.Data.OptStatus {
	case proto.OPT_STATUS_OFF:
		//  拉闸(关闭)状态
		return s.doStatusOff(ctx, md, vc)
	case proto.OPT_STATUS_ON:
		//  合闸(打开)状态
		return s.doStatusOn(ctx, md, vc)
	}

	// 程序不应该运行到这里
	s.Logger.ErrorContext(ctx, "unknow status", slog.Any("meter", md))
	return nil
}
