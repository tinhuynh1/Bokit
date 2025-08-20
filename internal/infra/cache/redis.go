package cache

import (
	"booking-svc/config"
	"context"
	"encoding/base64"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	password, err := base64.StdEncoding.DecodeString(cfg.Password)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisAddr(),
		Password: string(password),
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
