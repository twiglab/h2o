package chrgg

import (
	"context"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/chrgg/cdr"
)

func Handle(s *ChangeServer) mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		var cd cdr.ChargeData
		if err := cd.UnmarshalBinary(msg.Payload()); err != nil {
			log.Print(err)
			return
		}

		if _, err := s.DoChange(context.Background(), cd); err != nil {
			log.Print(err)
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
	// 连接到 Broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}
