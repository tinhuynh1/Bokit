package message_broker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	nc *nats.Conn
}

var PublisherConn *Publisher

type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	Source    string      `json:"source"`
}

func NewPublisher(url string) (*Publisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	PublisherConn = &Publisher{nc: nc}
	return PublisherConn, nil
}

func (p *Publisher) Publish(subject string, messageType string, data interface{}) error {
	message := Message{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
		Source:    "booking-service",
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.nc.Publish(subject, jsonData)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func ClosePublisher() {
	if PublisherConn != nil {
		PublisherConn.nc.Close()
	}
}
