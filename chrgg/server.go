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

func (s *ChargeServer) statusOff(ctx context.Context, md Meter, vcc *ent.VVC) error {
	//  拉闸状态，断开
	if md.Data.DataValue > vcc.Quota {
		// 大于限额，未触发, 保持状态
		return nil
	}

	// 小于限额，发送合闸消息，开
	ot := OnOffMsg{
		OP:      "ON",
		Device:  md.Device,
		Gateway: md.Gateway,
	}

	return s.Sender.SendData(ctx, ot)
}

func (s *ChargeServer) statusOn(ctx context.Context, md Meter, vcc *ent.VVC) error {
	//  合闸状态, 连通
	if md.Data.DataValue <= vcc.Quota {
		// 小于限额，未触发, 保持状态
		return nil
	}
	// 大于限额，发送拉闸消息，关
	ot := OnOffMsg{
		OP:      "OFF",
		Device:  md.Device,
		Gateway: md.Gateway,
	}

	return s.Sender.SendData(ctx, ot)
}

func (s *ChargeServer) Charge(ctx context.Context, md Meter) error {
	c, notfount, err := s.DBx.LoadLast(ctx, md.Code, md.Type)
	if notfount {
		return err
	}

	switch md.Data.OptStatus {
	case common.OPT_STATUS_OFF:
		//  拉闸状态
		return s.statusOff(ctx, md, c)
	case common.OPT_STATUS_ON:
		//  合闸状态
		return s.statusOn(ctx, md, c)
	}

	// 状态未知，无法处理
	return nil
}
