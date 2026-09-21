package chrgg

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const CLIENT_ID = "chrgg"

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

type MQTTAction struct {
	client mqtt.Client
}

func NewMQTTAction(client mqtt.Client) *MQTTAction {
	return &MQTTAction{client: client}
}

func (c *MQTTAction) SendData(ctx context.Context, obj SendObject) error {
	bb, err := obj.MarshalBinary()
	if err != nil {
		return err
	}

	pubToken := c.client.Publish(obj.Topic(), 0x01, false, bb)
	pubToken.Wait()

	return pubToken.Error()
}
