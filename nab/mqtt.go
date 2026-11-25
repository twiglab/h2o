package nab

import (
	"context"
	"errors"
	"log/slog"
	"time"

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
	go func() {
		// 异步非阻塞处理返回结果
		// 对于pub失败情况直接记录日志
		// 可观测性问题留到v2版本解决
		// TODO
		if err := mqtt.WaitTokenTimeout(token, 5*time.Second); err != nil {
			slog.ErrorContext(ctx, "mqtt Publish 失败", slog.Any("error", err))
			return
		}
	}()
	return nil
}

func ClientID(code, typ string) string {
	return code + "@" + typ
}
