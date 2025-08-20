package repository

import (
	"booking-svc/config"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(c *redis.Client) *RedisRepo {
	return &RedisRepo{client: c}
}

func (r *RedisRepo) StoreToken(ctx context.Context, tokenID string, ttl time.Duration) error {
	return r.client.Set(ctx, config.RedisTokenPrefix+tokenID, "valid", ttl).Err()
}

func (r *RedisRepo) IsTokenValid(ctx context.Context, tokenID string) bool {
	val, err := r.client.Get(ctx, config.RedisTokenPrefix+tokenID).Result()
	return err == nil && val == "valid"
}

func (r *RedisRepo) DeleteToken(ctx context.Context, tokenID string) error {
	return r.client.Del(ctx, config.RedisTokenPrefix+tokenID).Err()
}
