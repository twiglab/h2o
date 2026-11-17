package vigil

import (
	"context"
	"fmt"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const APP_ID = "vigil"

func Handle(s *Hub) mqtt.MessageHandler {
	ctx := context.Background()
	if s.BaseContext != nil {
		ctx = s.BaseContext(s)
	}

	return func(cli mqtt.Client, msg mqtt.Message) {
		var wm Meter
		if err := wm.UnmarshalBinary(msg.Payload()); err != nil {
			s.Logger.ErrorContext(ctx, "unmarshal error", slog.Any("err", err))
			return
		}
		if err := s.Handle(ctx, wm); err != nil {
			s.Logger.ErrorContext(ctx, "handle error", slog.Any("err", err))
			return
		}
	}
}

func RawHandle() mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		fmt.Println(string(msg.Payload()))
		//slog.Info("raw", slog.Any("raw", md))
	}
}
