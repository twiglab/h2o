package vigil

import (
	"context"
	"encoding/json/v2"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/pkg/common"
)

const CLIENT_ID = "vigil"

type Meter struct {
	common.Device
	Pos common.Pos `json:"pos,omitzero"`
}

type ElectricityMeter struct {
	Meter
	Data common.Electricity `json:"data,omitzero"`
}

func (d *ElectricityMeter) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func Handle(s *Hub) mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			return
		}

		switch common.TopicType(msg.Topic()) {
		case common.WaterTopic:
		case common.ElectricityTopic:
			var em ElectricityMeter
			if err := em.UnmarshalBinary(msg.Payload()); err != nil {
			}
			if err := s.HandleElectricity(context.Background(), em); err != nil {
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
