package service

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/jeremy2566/octopipe/internal/dao"
	"github.com/jeremy2566/octopipe/internal/model"
	"go.uber.org/zap"
	"resty.dev/v3"
)

var _ Zadig = &zadigImpl{}

type Zadig interface {
	GetTestEnvList(projectKey string) ([]model.RespZadigEnv, error)
	GetTestEnvDetail(envKey, projectKey string) (*model.RespZadigEnvDetail, error)
	Allocator(req model.AllocatorReq) error
	CreateSubEnv() error
	GetServiceCharts() map[string]string
}

type zadigImpl struct {
	log    *zap.Logger
	client *resty.Client
	rdb    dao.Rdb
}

func NewZadig(log *zap.Logger, client *resty.Client) Zadig {
	client.SetBaseURL("https://zadigx.shub.us").
		SetAuthToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiUnVpLkppYW5nIiwiZW1haWwiOiJydWkuamlhbmdAc3RvcmVodWIuY29tIiwidWlkIjoiMzBjYmZiZTAtNmYyNi0xMWVmLWEwYzEtNDI0Y2Q2NGY0MTZhIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiRGVyYWl2ZW4iLCJmZWRlcmF0ZWRfY2xhaW1zIjp7ImNvbm5lY3Rvcl9pZCI6ImdpdGh1YiIsInVzZXJfaWQiOiJEZXJhaXZlbiJ9LCJhdWQiOiJ6YWRpZyIsImV4cCI6NDg3OTUzOTU5Nn0.28147NOIPyGsFfuasHwHJlWvGAKSXCtn1oCD_J7vulM")
	return &zadigImpl{
		log:    log,
		client: client,
		rdb:    dao.NewRdb(log),
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

func (z *zadigImpl) Allocator(req model.AllocatorReq) error {
	switch {
	case req.BranchName == "main" || req.BranchName == "master":
		z.log.Info("master/main handle.")
	case req.BranchName == "release":
		z.log.Info("release handle.")
	case req.BranchName == "dev" || req.BranchName == "develop":
		z.log.Info("dev/develop handle.")
	case strings.HasPrefix(req.BranchName, "feat/") || strings.HasPrefix(req.BranchName, "feature/"):
		return z.handleFeature(req)
	default:
		z.log.Warn("unknown branch name", zap.String("branch", req.BranchName))
		return fmt.Errorf("unknown branch name: %s", req.BranchName)
	}

	return nil
}

func (z *zadigImpl) handleFeature(req model.AllocatorReq) error {
	branchName := req.BranchName
	namespace, err := z.rdb.GetNamespaceByBranch(branchName)
	if err != nil {
		return err
	}
	if namespace == nil {
		// 创建 ns，并部署服务
		z.CreateSubEnv()
	} else {
		// 添加新服务到 ns
	}

	return nil
}

func (z *zadigImpl) CreateSubEnv() error {
	namespace := z.selectedNamespace()
	se := model.ShareEnvReq{
		Enable:  true,
		IsBase:  false,
		BaseEnv: "test33",
	}
	req := model.CreateSubEnvReq{
		{
			EnvName:    namespace,
			ClusterID:  "64e48b5b8fc410571753cc6c",
			RegistryID: "64e485c78fc410571753cc67",
			Namespace:  namespace,
			IsExisted:  false,
			ShareEnv:   se,
		},
	}

	apiResponse := struct {
		Message string `json:"message"`
	}{}

	resp, err := z.client.R().
		SetResult(&apiResponse).
		SetContentType("application/json").
		SetBody(req).
		Post("/api/aslan/environment/environments?type=helm&projectName=fat-base-envrionment")
	if err != nil {
		z.log.Warn("failed to fetch environments", zap.Error(err))
	}
	if resp.StatusCode() != http.StatusOK {
		z.log.Warn("resp status code not ok",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response body", resp.String()),
		)
	}
	return nil
}

func (z *zadigImpl) selectedNamespace() string {
	envs, err := z.GetTestEnvList("fat-base-envrionment")
	if err != nil {
		z.log.Error("get zadig env list err", zap.Error(err))
		return ""
	}

	usedNamespaces := make(map[string]bool)
	for _, e := range envs {
		usedNamespaces[e.Namespace] = true
	}

	z.log.Info("fetch environments success", zap.Any("usedNamespaces", usedNamespaces))

	// 存储所有可用的命名空间
	var availableNamespaces []string
	for i := 1; i <= 50; i++ {
		ns := fmt.Sprintf("test%d", i)
		if !usedNamespaces[ns] {
			availableNamespaces = append(availableNamespaces, ns)
		}
	}

	// 检查是否有可用的命名空间
	if len(availableNamespaces) == 0 {
		z.log.Warn("No available namespaces in the 1-50 pool.")
	}
	randomIndex := rand.Intn(len(availableNamespaces))
	return availableNamespaces[randomIndex]
}

func (z *zadigImpl) GetServiceCharts() map[string]string {
	apiResponse := struct {
		ChartInfos []model.ChartInfoReq `json:"chart_infos"`
	}{}
	resp, err := z.client.R().
		SetResult(&apiResponse).
		Get("/api/aslan/environment/init_info/fat-base-envrionment?envType=share&isBaseEnv=false&baseEnv=test33&projectName=fat-base-envrionment")
	if err != nil {
		z.log.Warn("failed to fetch environments", zap.Error(err))
	}
	if resp.StatusCode() != http.StatusOK {
		z.log.Warn("resp status code not ok",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response body", resp.String()),
		)
	}
	ret := make(map[string]string)
	for _, info := range apiResponse.ChartInfos {
		ret[info.ServiceName] = info.ChartVersion
	}
	return ret
}
