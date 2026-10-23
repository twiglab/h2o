package mdna

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/proto"
)

type Server struct {
	Speaker Speaker
	MCli    mqtt.Client
	DBx     *DBx

	Logger *slog.Logger

	context context.Context
}

func (s *Server) Run(ctx context.Context) error {
	s.context = ctx

	msgMap := map[string]byte{
		"h2o/chrgg/#": 0x01,
		"h2o/opt/#":   0x01,
	}

	t := s.MCli.SubscribeMultiple(msgMap, s.msgHandle())
	t.Wait()
	return t.Error()
}

func (s *Server) msgHandle() mqtt.MessageHandler {

	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			msg.Ack()
			return
		}

		defer msg.Ack()

		parts := proto.TopicPart(msg.Topic())
		topicType := parts[1]

		switch topicType {
		case proto.QUOTA: // chrgg 发送的限额消息
			var chg ChrggMessage
			if err := chg.UnmarshalBinary(msg.Payload()); err != nil {
				// log
				return
			}
			if err := s.OnChrgg(s.context, chg); err != nil {
				// log
				return
			}
		case proto.OPT: // 盒子发送的开关消息(optStatus变更)
			code, op := parts[3], parts[4]
			if err := s.OnOpt(s.context, code, op); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func (s *Server) OnOpt(ctx context.Context, code, op string) error {
	if op == proto.ON {
		return nil
	}

	dev, err := s.DBx.LoadDeviceByCode(ctx, code)
	if err != nil {
		return err
	}

	n := StatusChangeNoticeBody{
		Now:    time.Now(),
		Code:   code,
		Op:     op,
		Device: dev,
	}
	return s.Speaker.StatusChangeNotice(ctx, n)
}

func (s *Server) OnChrgg(ctx context.Context, chg ChrggMessage) error {
	return nil
}
