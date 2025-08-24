package message_broker

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

var ConsumerConn *Consumer

type Consumer struct {
	nc *nats.Conn
}

func NewConsumer(url string) (*Consumer, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	return &Consumer{nc: nc}, nil
}

func (c *Consumer) Subscribe(subject string, callback func(msg *nats.Msg)) error {
	_, err := c.nc.Subscribe(subject, func(msg *nats.Msg) {
		callback(msg)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject: %w", err)
	}
	return nil
}

func CloseConsumer() {
	if ConsumerConn != nil {
		ConsumerConn.nc.Close()
	}
}
