package handler

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"resty.dev/v3"
)

type Handler struct {
	log    *zap.Logger
	client *resty.Client
	rdb    *redis.Client
}

func New(log *zap.Logger) *Handler {
	client := resty.New()
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(5 * time.Second)

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PWD"),
	})
	result, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Error("redis connect failed.", zap.Error(err))
	}
	log.Info("successfully connected to Redis.", zap.String("result", result))
	return &Handler{
		log:    log,
		client: client,
		rdb:    rdb,
	}
}
