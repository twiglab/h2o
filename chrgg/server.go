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

func isAlarm3(vc *ent.Top) bool {
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
				slog.Int64("top", vc.Top),
				slog.Int64("dataValue", md.Data.DataValue),
				slog.Int("status", vc.Status),
				slog.Any("meter", md), slog.Any("topRec", vc))
		}
	}

	if s.Alarm != 0 {
		if (vc.Top - md.Data.DataValue) < s.Alarm {
			if !isAlarm3(vc) { // 没拉闸报警过
				// 拉闸报警一次
				_ = vc.Update().
					SetAlarm(1).
					SetAlarmStock(md.Data.DataValue).
					SetAlarmTime(time.Now()).
					Exec(ctx) // 设置报警状态

				ot := newOnOffMessage(md, vc, common.OFF)
				return s.Sender.SendData(ctx, ot)
			}
		}
	}

	return nil
}

func (s *ChargeServer) Charge(ctx context.Context, md Meter) error {
	if md.Data.OptStatus == common.OPT_STATUS_UNKNOW {
		s.Logger.DebugContext(ctx, "status unknow", slog.Any("meter", md))
		return nil
	}

	c, notfount, err := s.DBx.LoadLast(ctx, md.Code, md.Type)
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
	// 增加这个字段的意义就是在不改变top的情况下，不计费
	// 这里有个问题要注意，找个状态是记录在当前限额记录上的，记录必须有效
	// 后续这个状态会移除，仅限当前版本使用
	if c.Status < STATUS_BEGIN {
		s.Logger.DebugContext(ctx, "人为指定 status < 0 强制不计费", slog.Any("meter", md), slog.Int("status", c.Status))
		return nil
	}

	switch md.Data.OptStatus {
	case common.OPT_STATUS_OFF:
		//  拉闸(关闭)状态
		return s.doStatusOff(ctx, md, c)
	case common.OPT_STATUS_ON:
		//  合闸(打开)状态
		return s.doStatusOn(ctx, md, c)
	}

	// 程序不应该运行到这里
	s.Logger.ErrorContext(ctx, "unknow status", slog.Any("meter", md))
	return nil
}
