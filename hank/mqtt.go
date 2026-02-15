package hank

import (
	"context"
	"encoding"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTAction struct {
	client mqtt.Client
}

func NewMQTTAction(client mqtt.Client) *MQTTAction {
	return &MQTTAction{client: client}
}

func (c *MQTTAction) SendData(ctx context.Context, data encoding.BinaryMarshaler) error {
	bb, err := data.MarshalBinary()
	if err != nil {
		return err
	}

	topic := ""

	pubToken := c.client.Publish(topic, 0x00, false, bb)
	pubToken.Wait()

	return pubToken.Error()
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

func pushTopic(code, typ string) string {
	return "h2o/" + code + "/" + typ
}
