package chrgg

import (
	"context"
	"encoding/json/v2"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/pkg/common"
)

type MeterData struct {
	common.Device
	Pos  common.Pos        `json:"pos,omitzero"`
	Data common.MeterValue `json:"data"`
	Flag common.Flag       `json:"flag,omitzero"`
}

func (d *MeterData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func HandleChange(s *ChargeServer) mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			return
		}

		var md MeterData
		if err := md.UnmarshalBinary(msg.Payload()); err != nil {
			slog.Error("unmarshal error", slog.Any("error", err))
			return
		}

		if _, err := s.Charge(context.Background(), md); err != nil {
			slog.Error("change error", slog.Any("error", err))
		}
	}
}

func RawHandle() mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			return
		}

		var cd ChargeData
		if err := cd.UnmarshalBinary(msg.Payload()); err != nil {
			slog.Error("unmarshal error", slog.Any("error", err))
			return
		}
		slog.Info("raw", slog.Any("chargeDate", cd))
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
