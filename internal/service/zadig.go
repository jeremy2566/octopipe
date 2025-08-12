package service

import (
	"fmt"
	"net/http"

	"github.com/jeremy2566/octopipe/internal/model"
	"go.uber.org/zap"
	"resty.dev/v3"
)

var _ Zadig = &zadigImpl{}

type Zadig interface {
	GetTestEnvList(projectKey string) ([]model.RespZadigEnv, error)
	GetTestEnvDetail(envKey, projectKey string) (*model.RespZadigEnvDetail, error)
}

type zadigImpl struct {
	log    *zap.Logger
	client *resty.Client
}

func NewZadig(log *zap.Logger, client *resty.Client) Zadig {
	client.SetBaseURL("https://zadigx.shub.us").
		SetAuthToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiUnVpLkppYW5nIiwiZW1haWwiOiJydWkuamlhbmdAc3RvcmVodWIuY29tIiwidWlkIjoiMzBjYmZiZTAtNmYyNi0xMWVmLWEwYzEtNDI0Y2Q2NGY0MTZhIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiRGVyYWl2ZW4iLCJmZWRlcmF0ZWRfY2xhaW1zIjp7ImNvbm5lY3Rvcl9pZCI6ImdpdGh1YiIsInVzZXJfaWQiOiJEZXJhaXZlbiJ9LCJhdWQiOiJ6YWRpZyIsImV4cCI6NDg3OTUzOTU5Nn0.28147NOIPyGsFfuasHwHJlWvGAKSXCtn1oCD_J7vulM")
	return &zadigImpl{
		log:    log,
		client: client,
	}
}

func (z *zadigImpl) GetTestEnvList(projectKey string) ([]model.RespZadigEnv, error) {
	var envs []model.RespZadigEnv
	resp, err := z.client.R().
		SetResult(&envs).
		SetQueryParam("projectKey", projectKey).
		Get("openapi/environments")
	if err != nil {
		z.log.Error("API 请求失败", zap.Error(err))
		return nil, fmt.Errorf("API request error: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		z.log.Error("API 请求返回非 OK 状态",
			zap.String("status", resp.Status()),
			zap.String("body", resp.String()))
		return nil, fmt.Errorf("API call failed with status: %s", resp.Status())
	}

	z.log.Info("GetTestEnvList", zap.Int("env length(include test17 and test33", len(envs)), zap.Any("envs", envs))
	return envs, nil
}

func (z *zadigImpl) GetTestEnvDetail(envKey, projectKey string) (*model.RespZadigEnvDetail, error) {
	var ret model.RespZadigEnvDetail

	resp, err := z.client.R().SetResult(&ret).
		SetPathParam("env_key", envKey).
		SetQueryParam("projectKey", projectKey).
		Get("/openapi/environments/{env_key}")
	if err != nil {
		z.log.Error("API 请求失败", zap.Error(err))
		return nil, fmt.Errorf("API request error: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		z.log.Error("API 请求返回非 OK 状态",
			zap.String("status", resp.Status()),
			zap.String("body", resp.String()))
		return nil, fmt.Errorf("API call failed with status: %s", resp.Status())
	}
	z.log.Info("GetTestEnvDetail", zap.Any("ret", ret))
	return &ret, nil
}
