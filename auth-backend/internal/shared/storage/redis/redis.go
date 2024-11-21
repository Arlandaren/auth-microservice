package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"service/internal/shared/config"
)

func ConnectRedis(ctx context.Context) *redis.Client {
	cfg := config.GetRedis()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}

	return rdb
}
