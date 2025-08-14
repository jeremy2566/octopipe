package service

import (
	"fmt"

	"github.com/jeremy2566/octopipe/internal/dao"
	"github.com/jeremy2566/octopipe/internal/model"
	"go.uber.org/zap"
	"resty.dev/v3"
)

var _ Cache = &cacheImpl{}

type Cache interface {
	ViewAll() []model.DaoNamespace
	DeleteNamespace(subEnv string) error
	SyncCache(projectKey string) error
	ViewBySubEnv(subEnv string) (*model.DaoNamespace, error)
}

type cacheImpl struct {
	redisDao     dao.Rdb
	log          *zap.Logger
	zadigService Zadig
}

func NewCache(
	log *zap.Logger,
	client *resty.Client,
) Cache {
	return &cacheImpl{
		redisDao:     dao.NewRdb(log),
		log:          log,
		zadigService: NewZadig(log, client),
	}
}

func (c *cacheImpl) SyncCache(projectKey string) error {
	c.log.Info("SyncCache")
	envs, err := c.zadigService.GetTestEnvList(projectKey)
	if err != nil {
		return fmt.Errorf("get test env list err: %w", err)
	}
	for _, env := range envs {
		if env.HasTest17OrTest33() {
			continue
		}
		detail, err := c.zadigService.GetTestEnvDetail(env.EnvKey, projectKey)
		if err != nil {
			c.log.Warn("get test env detail err", zap.String("sub env", env.Namespace), zap.Error(err))
			continue
		}

		cache, _ := c.redisDao.GetValueByKey(env.Namespace)
		var bn string
		if cache == nil {
			bn = "in-tree"
		} else {
			bn = cache.Branch
		}
		err = c.redisDao.SaveNamespace(env.Namespace, model.DaoNamespace{
			SubEnv:      env.Namespace,
			UpdateBy:    env.UpdateBy,
			Branch:      bn,
			ServiceName: detail.GetServices(),
		})
		if err != nil {
			c.log.Warn("save ns failed.", zap.String("namespace", env.Namespace), zap.Error(err))
			continue
		}
	}
	return nil
}

func (c *cacheImpl) ViewAll() []model.DaoNamespace {
	return c.redisDao.GetAllNamespace()
}

func (c *cacheImpl) ViewBySubEnv(subEnv string) (*model.DaoNamespace, error) {
	if err := c.SyncCache("fat-base-envrionment"); err != nil {
		return nil, err
	}

	namespaces := c.ViewAll()
	for _, namespace := range namespaces {
		if namespace.SubEnv == subEnv {
			return &namespace, nil
		}
	}
	return nil, fmt.Errorf("not found sub env: %s", subEnv)
}

func (c *cacheImpl) DeleteNamespace(subEnv string) error {
	env, err := c.ViewBySubEnv(subEnv)
	if err != nil {
		return fmt.Errorf("selected %w failed", err)
	}
	if env != nil {
		return fmt.Errorf("sub env[%s] has existed in zadig. Don't allow delete", subEnv)
	}
	return c.redisDao.DeleteNamespace(subEnv)
}
