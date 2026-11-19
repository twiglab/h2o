package nab

import (
	"context"
	"errors"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const NAB = "nab"

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

	if !c.client.IsConnectionOpen() {
		slog.ErrorContext(ctx, "mqtt not connected",
			slog.String("topic", obj.Topic()))
		return errors.New("mqtt not connected") // 必须返回错误，不能返回 nil
	}

	token := c.client.Publish(obj.Topic(), 0x00, false, bb)
	return token.Error()
}

func ClientID(code string) string {
	return code + "@" + NAB
}
