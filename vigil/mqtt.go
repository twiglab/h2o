package vigil

import (
	"context"
	"fmt"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const CLIENT_ID = "vigil"

func Handle(s *Hub) mqtt.MessageHandler {
	ctx := context.Background()
	if s.BaseContext != nil {
		ctx = s.BaseContext(s)
	}

	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			msg.Ack()
			return
		}

		defer msg.Ack()

		var wm Meter
		if err := wm.UnmarshalBinary(msg.Payload()); err != nil {
			s.Logger.ErrorContext(ctx, "unmarshal water error", slog.Any("err", err))
			return
		}
		if err := s.Handle(ctx, wm); err != nil {
			s.Logger.ErrorContext(ctx, "handle water error", slog.Any("err", err))
			return
		}

	}
}

func NewMQTTClient(clientID string, broker string, others ...string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.SetClientID(clientID)

	opts.AddBroker(broker)
	for _, b := range others {
		opts.AddBroker(b)
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

func RawHandle() mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			msg.Ack()
			return
		}
		defer msg.Ack()

		fmt.Println(string(msg.Payload()))
		//slog.Info("raw", slog.Any("raw", md))
	}
}
