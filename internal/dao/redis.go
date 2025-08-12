package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/jeremy2566/octopipe/internal/model"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var _ Rdb = &redisImpl{}

type Rdb interface {
	SaveNamespace(key string, ns model.DaoNamespace) error
	GetNamespace(key string) (*model.DaoNamespace, error)
	GetAllNamespace() []model.DaoNamespace
	DeleteNamespace(key string) error
}

type redisImpl struct {
	log    *zap.Logger
	client *redis.Client
}

func (r *redisImpl) DeleteNamespace(key string) error {
	//TODO implement me
	panic("implement me")
}

func NewRdb(log *zap.Logger) Rdb {
	return &redisImpl{
		log: log,
		client: redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: "redishello521@",
		}),
	}
}

func (r *redisImpl) GetAllNamespace() []model.DaoNamespace {
	var namespaces []model.DaoNamespace
	ctx := context.Background()

	// Use SCAN to safely iterate over keys in production without blocking.
	// The pattern "octopipe:namespace:*" will find all relevant keys.
	iter := r.client.Scan(ctx, 0, r.namespaceKey("*"), 0).Iterator()
	for iter.Next(ctx) {
		redisKey := iter.Val()
		data, err := r.client.HGetAll(ctx, redisKey).Result()
		if err != nil {
			r.log.Warn("failed to HGetAll for a key during GetAllNamespace", zap.String("key", redisKey), zap.Error(err))
			continue
		}

		if len(data) == 0 {
			continue
		}

		var serviceNames []string
		if err := json.Unmarshal([]byte(data["service_name"]), &serviceNames); err != nil {
			r.log.Warn("failed to unmarshal service_name for a key during GetAllNamespace", zap.String("key", redisKey), zap.Error(err))
			continue
		}

		namespaces = append(namespaces, model.DaoNamespace{
			SubEnv:      data["sub_env"],
			UpdateBy:    data["update_by"],
			Branch:      data["branch"],
			ServiceName: serviceNames,
		})
	}

	// Check for errors that occurred during the scan
	if err := iter.Err(); err != nil {
		r.log.Error("error during Rdb SCAN in GetAllNamespace", zap.Error(err))
	}

	return namespaces
}

func (r *redisImpl) GetNamespace(key string) (*model.DaoNamespace, error) {
	redisKey := r.namespaceKey(key)
	data, err := r.client.HGetAll(context.Background(), redisKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to HGetAll from Rdb: %w", err)
	}

	if len(data) == 0 {
		return nil, nil
	}

	var serviceNames []string
	if err := json.Unmarshal([]byte(data["service_name"]), &serviceNames); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service_name from Rdb: %w", err)
	}

	ns := &model.DaoNamespace{
		SubEnv:      key,
		Branch:      data["branch"],
		ServiceName: serviceNames,
	}

	return ns, nil
}

func (r *redisImpl) SaveNamespace(key string, ns model.DaoNamespace) error {
	redisKey := r.namespaceKey(key)
	r.log.Info("Attempting to HSET to Rdb",
		zap.String("redis_key", redisKey),
		zap.String("update_by", ns.UpdateBy),
		zap.String("branch", ns.Branch),
		zap.Any("service_name", ns.ServiceName))
	services, _ := json.Marshal(ns.ServiceName)
	err := r.client.HSet(context.Background(), redisKey,
		"sub_env", ns.SubEnv,
		"update_by", ns.UpdateBy,
		"branch", ns.Branch,
		"service_name", services,
	).Err()
	if err != nil {
		return fmt.Errorf("failed to HSET to Rdb: %w", err)
	}
	return nil
}

// namespaceKey is a helper function to generate a consistent Rdb key for a namespace.
func (r *redisImpl) namespaceKey(key string) string {
	return fmt.Sprintf("octopipe:namespace:%s", key)
}

func NewRedis(log *zap.Logger, client *redis.Client) Rdb {
	return &redisImpl{
		log:    log,
		client: client,
	}
}
