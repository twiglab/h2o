package mdna

import (
	"context"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/dbx"
	"github.com/twiglab/h2o/pkg/common"
)

type Server struct {
	DBx  *dbx.DBx
	MCli mqtt.Client

	Speaker Speaker

	Logger *slog.Logger

	Context context.Context

	Alarm1 int64
	Alarm2 int64
	Alarm3 int64
}

func (s *Server) Run(ctx context.Context) error {
	s.Context = ctx

	msgMap := map[string]byte{
		"h2o/chrgg/#": 0x01,
		"h2o/opt/#":   0x01,
	}

	t := s.MCli.SubscribeMultiple(msgMap, s.MsgHandle())
	t.Wait()
	return t.Error()
}

func (s *Server) MsgHandle() mqtt.MessageHandler {

	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			return
		}

		defer msg.Ack()

		parts := common.TopicPart(msg.Topic())
		topicType := parts[1]

		switch topicType {
		case "chrgg": // chrgg 发送的限额消息
			var chg ChrggMessage
			if err := chg.UnmarshalBinary(msg.Payload()); err != nil {
				// log
				return
			}
			if err := s.OnChrgg(s.Context, chg); err != nil {
				// log
				return
			}
		case "opt": // 盒子发送的开关消息(optStatus变更)
			code, op := parts[3], parts[4]
			if err := s.OnOpt(s.Context, code, op); err != nil {
				// log
				return
			}
		}
	}
}

func (s *Server) OnOpt(ctx context.Context, code, op string) error {
	if op == common.ON {
		return nil
	}

	dev, err := s.DBx.LoadDevice(ctx, code)
	if err != nil {
		return err
	}

	n := Notice{
		Device: dev,
		Type:   "opt",
		Op:     op,
	}
	return s.Speaker.Notice(ctx, n)
}

func (s *Server) OnChrgg(ctx context.Context, chg ChrggMessage) error {

	if chg.Top.Status != 0 {
		// 已经结束或者非法的限额记录
		return nil
	}

	if chg.Data.DataValue > chg.Top.Top {
		// 人为调整的记录
		return nil
	}

	dev, err := s.DBx.LoadDevice(ctx, chg.Code)
	if err != nil {
		return nil
	}
	n := Notice{
		Device: dev,
		Charge: chg,
		Type:   "chrgg",
	}
	return s.Speaker.Notice(ctx, n)
}
