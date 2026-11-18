package nab

import (
	"context"
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

	if !c.client.IsConnected() {
		//  TODO
		// 此处逻辑有缺陷，记录日志，暂时返回nil
		slog.ErrorContext(ctx, "mqtt not connected")
		return nil
	}

	token := c.client.Publish(obj.Topic(), 0x00, false, bb)
	return token.Error()
}

func ClientID(code string) string {
	return code + "@" + NAB
}
