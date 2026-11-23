package nab

import (
	"context"
	"errors"
	"io"
	"log"
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

	// 允许断线时，数据丢失
	// 数据由采集轮询发送，采集为最小也有300毫秒左右，无背压风险
	if !c.client.IsConnectionOpen() {
		slog.ErrorContext(ctx, "mqtt not connected",
			slog.String("topic", obj.Topic()))
		return errors.New("mqtt not connected")
	}
	token := c.client.Publish(obj.Topic(), 0x00, false, bb)
	return token.Error()
}

func ClientID(code string) string {
	return code + "@" + NAB
}

func MqttLog(w io.Writer) *log.Logger {
	l := log.New(w, "mqtt", log.LstdFlags)
	return l
}
