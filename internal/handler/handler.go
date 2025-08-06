package handler

import (
	"go.uber.org/zap"
	"resty.dev/v3"
	"time"
)

type Handler struct {
	log    *zap.Logger
	client *resty.Client
}

func New(log *zap.Logger) *Handler {
	client := resty.New()
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(5 * time.Second)
	return &Handler{
		log:    log,
		client: client,
	}
}
