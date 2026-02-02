package hank

import (
	"bytes"
	"context"
	"encoding/json/v2"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/pkg/data"
)

type MQTTAction struct {
	client mqtt.Client
}

func NewMQTTAction(client mqtt.Client) *MQTTAction {
	return &MQTTAction{client: client}
}

func (c *MQTTAction) SendData(ctx context.Context, data data.Device) error {
	var bb bytes.Buffer
	bb.Grow(1024)

	err := json.MarshalWrite(&bb, &data)
	if err != nil {
		return err
	}

	topic := pushTopic(data.Code, data.Type)

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
	// 连接到 Broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

func pushTopic(code, typ string) string {
	return "h2o/" + code + "/" + typ
}
