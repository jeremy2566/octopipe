package api

import (
	"time"

	"github.com/jeremy2566/octopipe/internal/service"
	"go.uber.org/zap"
	"resty.dev/v3"
)

type Api struct {
	cache service.Cache
}

func New(log *zap.Logger) Api {
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)

	cache := service.NewCache(log, client)
	return Api{
		cache: cache,
	}
}
