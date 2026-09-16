package chrgg

import (
	"context"
	"log/slog"

	"github.com/twiglab/h2o/chrgg/orm"
	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/pkg/common"
)

type ChargeServer struct {
	DBx *orm.DBx

	Sender Sender

	Logger *slog.Logger
}

func (s *ChargeServer) OptOff(ctx context.Context, md Meter, vc *ent.ValueCharge) error {
	//  拉闸状态，断开
	if md.Data.DataValue > vc.Top {
		// 大于限额，未触发, 保持状态
		return nil
	}

	// 小于限额，发送合闸消息，开
	ot := OnOffMsg{
		OP:      common.ON,
		Device:  md.Device,
		Gateway: md.Gateway,
		Charge: Charge{
			Code:       vc.Code,
			Top:        vc.Top,
			Current:    md.Data.DataValue,
			Stock:      vc.Stock,
			Incr:       vc.Incr,
			Amount:     vc.Amount,
			UnitPrice:  vc.UnitPrice,
			ChargeTime: vc.ChargeTime,
		},
	}

	return s.Sender.SendData(ctx, ot)
}

func (s *ChargeServer) OptOn(ctx context.Context, md Meter, vc *ent.ValueCharge) error {
	//  合闸状态, 连通
	if md.Data.DataValue <= vc.Top {
		// 小于限额，未触发, 保持状态
		return nil
	}

	// 大于限额，发送拉闸消息，关
	ot := OnOffMsg{
		OP:      common.OFF,
		Device:  md.Device,
		Gateway: md.Gateway,
		Charge: Charge{
			Code:       vc.Code,
			Top:        vc.Top,
			Current:    md.Data.DataValue,
			Stock:      vc.Stock,
			Incr:       vc.Incr,
			Amount:     vc.Amount,
			UnitPrice:  vc.UnitPrice,
			ChargeTime: vc.ChargeTime,
		},
	}

	return s.Sender.SendData(ctx, ot)
}

func (s *ChargeServer) Charge(ctx context.Context, md Meter) error {
	c, notfount, err := s.DBx.LoadLast(ctx, md.Code, md.Type)
	if notfount {
		return err
	}

	// 其他错误
	if err != nil {
		return err
	}

	if c.Top < 0 {
		return nil
	}

	/*
		// 不使用status
		if c.Status < 0 {
			return nil
		}
	*/

	switch md.Data.OptStatus {
	case common.OPT_STATUS_OFF:
		//  拉闸(关闭)状态
		return s.OptOff(ctx, md, c)
	case common.OPT_STATUS_ON:
		//  合闸(打开)状态
		return s.OptOn(ctx, md, c)
	}

	// 状态未知，无法处理
	return nil
}
