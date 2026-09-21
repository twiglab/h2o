package chrgg

import (
	"cmp"
	"context"
	"log/slog"

	"github.com/twiglab/h2o/chrgg/orm"
	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/pkg/common"
)

const no_alarm = 0

func isAlarm3(vc *ent.ValueCharge) bool {
	return vc.Alarm3 != no_alarm
}

type ChargeServer struct {
	DBx *orm.DBx

	Sender Sender

	Logger *slog.Logger

	AlarmQuota int64
}

func (s *ChargeServer) OptOff(ctx context.Context, md Meter, vc *ent.ValueCharge) error {
	//  拉闸状态，断开
	if cmp.Less(md.Data.DataValue, vc.Top) {
		// 在断开状态，小于限额，发送合闸消息，开
		ot := NewOnOffMessage(md, vc, common.ON)
		return s.Sender.SendData(ctx, ot)
	}
	return nil
}

func (s *ChargeServer) OptOn(ctx context.Context, md Meter, vc *ent.ValueCharge) error {

	//  合闸状态, 连通
	if cmp.Less(vc.Top, md.Data.DataValue) {
		// 超额, 拉闸断开
		ot := NewOnOffMessage(md, vc, common.OFF)
		return s.Sender.SendData(ctx, ot)
	}

	if s.AlarmQuota != 0 {
		if (vc.Top - md.Data.DataValue) < s.AlarmQuota {
			if !isAlarm3(vc) { // 没拉闸报警过
				// 拉闸报警一次
				_ = vc.Update().SetAlarm3(1).Exec(ctx) // 设置报警状态
				ot := NewOnOffMessage(md, vc, common.OFF)
				return s.Sender.SendData(ctx, ot)
			}
		}
	}

	return nil
}

func (s *ChargeServer) Charge(ctx context.Context, md Meter) error {
	c, notfount, err := s.DBx.LoadLast(ctx, md.Code, md.Type)
	if notfount {
		return nil // 没找到限额，无法计费，默认不计费
	}

	// 其他错误
	if err != nil {
		// log
		return err
	}

	if c.Top < 0 {
		// log
		return nil // 强制不计费, 人为指定
	}

	// 当前限额的业务状态
	// 增加这个字段的意义就是在不改变top的情况下，不计费
	// 这里有个问题要注意，找个状态是记录在当前限额记录上的，记录必须有效
	// 后续这个状态会移除，仅限当前版本使用
	if c.Status < 0 {
		// 小于零，人为指定不计费
		// log
		return nil
	}

	switch md.Data.OptStatus {
	case common.OPT_STATUS_OFF:
		//  拉闸(关闭)状态
		return s.OptOff(ctx, md, c)
	case common.OPT_STATUS_ON:
		//  合闸(打开)状态
		return s.OptOn(ctx, md, c)
	}

	// 状态未知，无法处理
	// log
	return nil
}
