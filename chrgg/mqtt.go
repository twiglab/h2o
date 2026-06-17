package chrgg

import (
	"context"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/pkg/common"
)

const CLIENT_ID = "chrgg"

func HandleChange(s *ChargeServer) mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			return
		}

		defer msg.Ack()

		switch common.TopicType(msg.Topic()) {
		case common.WaterTopic:

		case common.ElectricityTopic:
			var em ElectyMeterData
			if err := em.UnmarshalBinary(msg.Payload()); err != nil {
				s.Logger.Error("unmarshal error", slog.Any("error", err))
				return
			}
			if _, err := s.Charge(context.Background(), em); err != nil {
				s.Logger.Error("charge error", slog.Any("raw", em), slog.Any("error", err))
			}
		case common.GasTopic:
		}
	}
}

func RawHandle() mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			return
		}

		var md ElectyMeterData
		if err := md.UnmarshalBinary(msg.Payload()); err != nil {
			slog.Error("unmarshal error", slog.Any("error", err))
			return
		}
		slog.Info("raw", slog.Any("raw", md))
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
	// 连接到 Broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}
