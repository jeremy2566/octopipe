package api

import (
	"time"

	"github.com/jeremy2566/octopipe/internal/service"
	"go.uber.org/zap"
	"resty.dev/v3"
)

type Api struct {
	cache service.Cache
	zadig service.Zadig
}

func New(log *zap.Logger) Api {
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)

	return Api{
		cache: service.NewCache(log, client),
		zadig: service.NewZadig(log, client),
	}
}
