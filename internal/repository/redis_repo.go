package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(c *redis.Client) *RedisRepo {
	return &RedisRepo{client: c}
}

func (r *RedisRepo) Publish(ctx context.Context, channel string, event interface{}) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = r.client.Publish(ctx, channel, eventData).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisRepo) Subscribe(ctx context.Context, channel string) (message string, err error) {
	pubsub := r.client.Subscribe(ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case msg := <-ch:
			if msg != nil {
				return msg.Payload, nil
			}
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}
