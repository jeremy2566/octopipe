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
		err = c.redisDao.SaveNamespace(env.Namespace, model.DaoNamespace{
			SubEnv:      env.Namespace,
			UpdateBy:    env.UpdateBy,
			Branch:      "in-tree",
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

func (c *cacheImpl) DeleteNamespace(subEnv string) error {
	return c.redisDao.DeleteNamespace(subEnv)
}
