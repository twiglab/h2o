package chrgg

import (
	"context"
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/pkg/common"
)

type RawData struct {
	common.Device
	Pos  common.Pos   `json:"pos,omitzero"`
	Data common.DataV `json:"data"`
	Flag common.Flag  `json:"flag,omitzero"`
}

func (d *RawData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func HandleChange(s *ChangeServer) mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			return
		}
		var rd RawData
		if err := rd.UnmarshalBinary(msg.Payload()); err != nil {
			log.Print(err)
			return
		}

		if _, err := s.DoChange(context.Background(), rd); err != nil {
			log.Print(err)
		}
	}
}

func RawHandle() mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		var cd ChargeData
		if err := cd.UnmarshalBinary(msg.Payload()); err != nil {
			log.Print(err)
			return
		}
		log.Print(cd)
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
