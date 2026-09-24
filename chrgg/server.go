package chrgg

import (
	"cmp"
	"context"
	"log/slog"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/chrgg/orm"
	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/clog/wal"
	"github.com/twiglab/h2o/pkg/common"
)

const no_alarm = 0

func isAlarm(vc *ent.Top) bool {
	return vc.Alarm != no_alarm
}

type ChargeServer struct {
	DBx    *orm.DBx
	Sender Sender

	Logger *slog.Logger
	WAL    *wal.WAL

	Alarm int64

	MCli mqtt.Client
}

func (s *ChargeServer) Run() error {
	t := s.MCli.Subscribe(common.GeneralDataTopic, 0x01, s.MsgHandle())
	t.Wait()
	return t.Error()
}

func (s *ChargeServer) MsgHandle() mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			return
		}

		defer msg.Ack()

		switch common.DataTopicType(msg.Topic()) {
		case common.GasDataTopic:
		case common.ElectricityDataTopic, common.WaterDataTopic: // 目前只支持水表和电表
			var em Meter
			if err := em.UnmarshalBinary(msg.Payload()); err != nil {
				s.Logger.Error("unmarshal error", slog.Any("error", err))
				return
			}
			if err := s.Charge(context.Background(), em); err != nil {
				s.Logger.Error("charge error", slog.Any("raw", em), slog.Any("error", err))
			}
		}
	}
}

func (s *ChargeServer) doStatusOff(ctx context.Context, md Meter, vc *ent.Top) error {
	//  拉闸状态，断开
	if cmp.Less(md.Data.DataValue, vc.Top) {
		// 在断开状态，小于限额，发送合闸消息，开
		s.WAL.WriteLogContext(ctx,
			wal.String("top.code", vc.Code),
			wal.Int64("top", vc.Top), wal.Int64("dataValue", md.Data.DataValue),
			wal.String("code", md.Code), wal.String("type", md.Type),
			wal.String("onoffType", "on"), // 超出限额，正常合闸操作
			wal.String("msg", "正常合闸"),
		)

		ot := newOnOffMessage(md, vc, common.ON)
		return s.Sender.SendData(ctx, ot)
	}
	return nil
}

func (s *ChargeServer) doStatusOn(ctx context.Context, md Meter, vc *ent.Top) error {
	//  合闸状态, 连通
	if cmp.Less(vc.Top, md.Data.DataValue) {
		// 超额, 拉闸断开
		if vc.Status == STATUS_BEGIN {
			// 限额记录正常，执行拉闸操作
			s.WAL.WriteLogContext(ctx,
				wal.String("top.code", vc.Code),
				wal.Int64("top", vc.Top), wal.Int64("dataValue", md.Data.DataValue),
				wal.String("code", md.Code), wal.String("type", md.Type),
				wal.String("onoffType", "off"), // 超出限额，正常拉闸操作
				wal.String("msg", "超出限额"),
			)

			_ = vc.Update().
				SetStatus(STATUS_END).
				SetEndTime(time.Now()).         // 计费结束时间
				SetEndStock(md.Data.DataValue). // 计费结束的表显
				Exec(ctx)                       // 计费结束

			// 断开
			ot := newOnOffMessage(md, vc, common.OFF)
			return s.Sender.SendData(ctx, ot)
		} else {
			// 非正常状态，疑似数据有非法修改
			s.Logger.WarnContext(ctx, "illegal status",
				slog.String("top.code", vc.Code),
				slog.Int64("top", vc.Top),
				slog.Int64("dataValue", md.Data.DataValue),
				slog.Int("status", vc.Status),
				slog.Any("meter", md), slog.Any("topRec", vc))
		}
	}

	// 以下是报警逻辑， 如果alarm 为 0 不报警
	if s.Alarm == 0 {
		return nil
	}

	if (vc.Top - md.Data.DataValue) < s.Alarm {
		if !isAlarm(vc) { // 没拉闸报警过
			// 拉闸报警一次
			s.WAL.WriteLogContext(ctx,
				wal.String("top.code", vc.Code),
				wal.Int64("top", vc.Top), wal.Int64("dataValue", md.Data.DataValue),
				wal.String("code", md.Code), wal.String("type", md.Type),
				wal.String("onoffType", "alarm"), // 超出限额，正常拉闸操作
				wal.String("msg", "报警拉闸"),
			)
			_ = vc.Update().
				SetAlarm(1).
				SetAlarmStock(md.Data.DataValue).
				SetAlarmTime(time.Now()).
				Exec(ctx) // 设置报警状态

			ot := newOnOffMessage(md, vc, common.OFF)
			return s.Sender.SendData(ctx, ot)
		}
	}

	return nil
}

func (s *ChargeServer) Charge(ctx context.Context, md Meter) error {
	if md.Data.OptStatus == common.OPT_STATUS_UNKNOW {
		s.Logger.DebugContext(ctx, "status unknow", slog.Any("meter", md))
		return nil
	}

	vc, notfount, err := s.DBx.LoadLast(ctx, md.Code, md.Type)
	if notfount {
		// 没找到限额，无法计费，默认不计费
		s.Logger.DebugContext(ctx, "record not found", slog.Any("meter", md))
		return nil
	}

	// 其他错误
	if err != nil {
		s.Logger.ErrorContext(ctx, "loadLast error", slog.Any("meter", md), slog.Any("error", err))
		return err
	}

	// 当前限额的业务状态
	// 这里有个问题要注意，找个状态是记录在当前限额记录上的，记录必须有效
	// 这里是最后一次人为阻止
	if vc.Status < STATUS_BEGIN {
		s.Logger.DebugContext(ctx, "人为指定 status < 0 强制不计费", slog.Any("meter", md), slog.Int("status", vc.Status))
		return nil
	}

	if vc.Status == STATUS_BEGIN {
		// 正常状态下，计费没有结束，发送chrgg消息
		m := ChrggMessage{
			Device:  md.Device,
			Pos:     md.Pos,
			Data:    md.Data,
			Gateway: md.Gateway,
			Top: Top{
				Code:       vc.Code,
				Top:        vc.Top,
				Current:    md.Data.DataValue,
				Stock:      vc.Stock,
				Incr:       vc.Incr,
				Amount:     vc.Amount,
				UnitPrice:  vc.UnitPrice,
				ChargeTime: vc.ChargeTime,

				Alarm:     vc.Alarm,
				AlarmTime: vc.AlarmTime,

				Status:  vc.Status,
				EndTime: vc.EndTime,
			},
		}
		s.Sender.SendData(ctx, m)
	}

	// 注意： 即便状态为 STATUS_END 也是要继续执行的
	// 是否开关，完全由限额决定
	switch md.Data.OptStatus {
	case common.OPT_STATUS_OFF:
		//  拉闸(关闭)状态
		return s.doStatusOff(ctx, md, vc)
	case common.OPT_STATUS_ON:
		//  合闸(打开)状态
		return s.doStatusOn(ctx, md, vc)
	}

	// 程序不应该运行到这里
	s.Logger.ErrorContext(ctx, "unknow status", slog.Any("meter", md))
	return nil
}
