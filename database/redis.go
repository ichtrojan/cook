package database

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/ichtrojan/cook/config"
	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

var ctx = context.Background()

func ConnectRedis(cfg config.Redis) error {
	if cfg.Password == "null" {
		cfg.Password = ""
	}

	var client *redis.Client

	if cfg.Scheme == "tls" {
		client = redis.NewClient(&redis.Options{
			Addr:       fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password:   cfg.Password,
			DB:         0,
			TLSConfig:  &tls.Config{},
			MaxRetries: 3,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:       fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password:   cfg.Password,
			DB:         0,
			MaxRetries: 3,
		})
	}

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	fmt.Println("connected to redis", config.RedisConfig)

	Redis = client

	return nil
}
