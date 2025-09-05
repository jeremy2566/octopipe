package service

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"resty.dev/v3"
)

func TestCacheImpl_ViewBySubEnv(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)
	cache := NewCache(
		log,
		client,
	)

	env, _ := cache.ViewBySubEnv("test8")
	log.Info("", zap.Any("env", env))
}

func TestCacheImpl_ViewAll(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)

	cache := NewCache(
		log,
		client,
	)
	all := cache.ViewAll()

	for _, namespace := range all {
		log.Info("", zap.Any("namespace", namespace))
	}
}

func TestCacheImpl_DeleteNamespace(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)
	cache := NewCache(
		log,
		client,
	)

	err := cache.DeleteNamespace("test8")
	log.Info("", zap.Error(err))
}
