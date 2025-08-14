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
	UpdateServiceByKey(key string, service string)
	GetValueByKey(key string) (*model.DaoNamespace, error)
	GetAllNamespace() []model.DaoNamespace
	DeleteNamespace(key string) error
	GetNamespaceByBranch(branch string) (*model.DaoNamespace, error)
}

type redisImpl struct {
	log    *zap.Logger
	client *redis.Client
}

func (r *redisImpl) UpdateServiceByKey(key string, service string) {
	ns, err := r.GetValueByKey(key)
	if err != nil {
		r.log.Warn("updateServiceByKey failed.", zap.Error(err))
		return
	}

	if ns == nil {
		r.log.Warn("namespace is nil.", zap.String("key", key))
		return
	}
	services := ns.ServiceName
	for _, s := range services {
		if s == service {
			return
		}
	}

	services = append(services, service)
	marshal, _ := json.Marshal(services)
	r.client.HSet(context.Background(), r.namespaceKey(key), "service_name", marshal)
}

func (r *redisImpl) DeleteNamespace(key string) error {
	redisKey := r.namespaceKey(key)
	r.log.Info("Deleting namespace from Redis", zap.String("key", redisKey))

	// DEL command returns the number of keys that were removed.
	// We just need to check for an error.
	err := r.client.Del(context.Background(), redisKey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s from Redis: %w", redisKey, err)
	}
	return nil
}

func NewRdb(log *zap.Logger) Rdb {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PWD"),
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis for DAO", zap.Error(err))
	}
	log.Info("Successfully connected to Redis for DAO")

	return &redisImpl{
		log:    log,
		client: client,
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

func (r *redisImpl) GetValueByKey(key string) (*model.DaoNamespace, error) {
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
	err := r.client.HSet(
		context.Background(), redisKey,
		"sub_env", key, // 修复：这里应该保存子环境名(key)，而不是完整的 Redis 键。
		"update_by", ns.UpdateBy,
		"branch", ns.Branch,
		"service_name", services,
	).Err()
	if err != nil {
		return fmt.Errorf("failed to HSET to Rdb: %w", err)
	}
	return nil
}

func (r *redisImpl) GetNamespaceByBranch(branch string) (*model.DaoNamespace, error) {
	// 注意：此实现会遍历所有命名空间。
	// 如果命名空间数量非常大，更高效的做法是建立一个从 branch 到 sub_env 的二级索引。
	allNamespaces := r.GetAllNamespace()

	for _, ns := range allNamespaces {
		if ns.Branch == branch {
			// 创建一个新的变量来持有找到的命名空间，并返回其地址。
			// 这可以避免返回一个指向循环变量 `ns` 的指针，那是一个常见的 Go 语言陷阱。
			result := ns
			return &result, nil
		}
	}

	// 没有找到匹配的 branch，这不是一个错误。
	return nil, nil
}

// namespaceKey is a helper function to generate a consistent Rdb key for a namespace.
func (r *redisImpl) namespaceKey(key string) string {
	return fmt.Sprintf("octopipe:namespace:%s", key)
}
