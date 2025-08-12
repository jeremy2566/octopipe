package dao

import (
	"os"
	"testing"

	"github.com/jeremy2566/octopipe/internal/model"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRedisImpl_SaveNamespace(t *testing.T) {
	log, _ := zap.NewDevelopment()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "redishello521@",
	})
	r := NewRedis(log, rdb)

	err := r.SaveNamespace("test-1", model.DaoNamespace{
		SubEnv:      "test-1",
		Branch:      "branch-1",
		ServiceName: []string{"service-1"},
	})
	assert.Nil(t, err)
}

func TestNewRedis_GetNamespace(t *testing.T) {
	log, _ := zap.NewDevelopment()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "redishello521@",
	})
	r := NewRedis(log, rdb)
	namespace, _ := r.GetNamespace("test-1")
	log.Info("namespace", zap.Any("namespace", namespace))
}

func TestNewRedis_GetAllNamespace(t *testing.T) {
	log, _ := zap.NewDevelopment()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "redishello521@",
	})
	r := NewRedis(log, rdb)
	namespace := r.GetAllNamespace()
	log.Info("namespace", zap.Any("ns", namespace))
}
