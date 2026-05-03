package hank

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsAction struct {
	js jetstream.JetStream
}

func (c *NatsAction) SendData(ctx context.Context, obj SendObject) error {
	bs, err := obj.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = c.js.Publish(ctx, obj.Topic(), bs)
	return err
}

func (c *NatsAction) Close() error {
	return c.js.Conn().Drain()
}

func NewNatsAction(url string) (*NatsAction, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	js, _ := jetstream.New(nc)
	return &NatsAction{
		js: js,
	}, nil
}
