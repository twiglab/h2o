package hank

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsAction struct {
	jetstream.Publisher
}

func (c *NatsAction) SendData(ctx context.Context, obj SendObject) error {
	bs, err := obj.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = c.Publish(ctx, obj.Topic(), bs)
	return err
}

func NewNatsAction() (*NatsAction, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	js, _ := jetstream.New(nc)
	return &NatsAction{
		Publisher: js,
	}, nil
}
