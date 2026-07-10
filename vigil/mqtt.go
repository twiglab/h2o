package vigil

import (
	"context"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/pkg/common"
)

const CLIENT_ID = "vigil"

func Handle(s *Hub) mqtt.MessageHandler {
	ctx := context.Background()
	if s.BaseContext != nil {
		ctx = s.BaseContext(s)
	}

	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			return
		}

		defer msg.Ack()

		switch common.TopicType(msg.Topic()) {
		case common.WaterTopic:
			var wm WaterMeter
			if err := wm.UnmarshalBinary(msg.Payload()); err != nil {
				s.Logger.ErrorContext(ctx, "unmarshal water error", slog.Any("err", err))
				return
			}
			wm.setup()
			if err := s.HandleWater(ctx, wm); err != nil {
				s.Logger.ErrorContext(ctx, "handle water error", slog.Any("err", err))
				return
			}

		case common.ElectricityTopic:
			var em ElectricityMeter
			if err := em.UnmarshalBinary(msg.Payload()); err != nil {
				s.Logger.ErrorContext(ctx, "unmarshal electy error", slog.Any("err", err))
				return
			}
			em.setup()
			if err := s.HandleElecty(ctx, em); err != nil {
				s.Logger.ErrorContext(ctx, "handle electy error", slog.Any("err", err))
				return
			}

		case common.GasTopic:
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
