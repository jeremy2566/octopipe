package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

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
	CreateSubEnv() (string, error)
	GetServiceCharts() map[string]string
	DeployService(req model.DeployServiceReq) (int, error)
	AddService(req model.AddServiceReq) error
	DeleteSubEnv(env string) error
}

type zadigImpl struct {
	log    *zap.Logger
	client *resty.Client
	rdb    dao.Rdb
	lark   Lark
}

func (z *zadigImpl) DeleteSubEnv(env string) error {
	return z.rdb.DeleteNamespace(env)
}

func (z *zadigImpl) AddService(req model.AddServiceReq) error {
	z.log.Info("add service", zap.Any("params", req))
	charts := z.GetServiceCharts()
	chartVersion, exist := charts[req.ServiceName]
	if !exist {
		z.log.Warn("service not found", zap.String("service", req.ServiceName))
		return fmt.Errorf("service not found")
	}

	ufcv := []model.UtilsFunChartValues{
		{
			EnvName:         req.SubEnv,
			ServiceName:     req.ServiceName,
			ReleaseName:     req.ServiceName,
			ChartVersion:    chartVersion,
			Deploy_strategy: "deploy",
		},
	}
	uf := model.UtilsFun{
		ReplacePolicy: "notUseEnvImage",
		EnvNames:      []string{req.SubEnv},
		ChartValues:   ufcv,
	}

	var apiResponse []struct {
		EnvName    string `json:"env_name"`
		Status     string `json:"status"`
		ErrMessage string `json:"err_message"`
	}
	resp, err := z.client.R().
		SetResult(&apiResponse).
		SetBody(uf).
		Put("/api/aslan/environment/environments?type=helm&projectName=fat-base-envrionment")

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

// trans 服务名转 service name, service module, repo name 和 build name
func (z *zadigImpl) trans(serviceName string) (string, string, string, string) {
	switch serviceName {
	case "backoffice-v1-web":
		return "backoffice-v1-web-app", "backoffice-v1-web", "backoffice-v1-web", "backoffice-v1-web"
	default:
		return serviceName, serviceName, serviceName, serviceName
	}
}

func (z *zadigImpl) DeployService(req model.DeployServiceReq) (int, error) {
	z.log.Info("deploy service.", zap.Any("params", req))
	sn, sm, rn, bn := z.trans(req.ServiceName)

	servicesReq := model.DeployServicesReq{
		Name:        "test33",
		DisplayName: "fat-base-workflow",
		Project:     "fat-base-envrionment",
		Params: []struct {
			Name  string `json:"name"`
			Type  string `json:"type"`
			Value string `json:"value"`
		}{
			{
				Name:  "环境",
				Type:  "choice",
				Value: req.SubEnv,
			},
		},
		Stages: []struct {
			Name string `json:"name"`
			Jobs []struct {
				Name string `json:"name"`
				Type string `json:"type"`
				Spec struct {
					DefaultServiceAndBuilds []struct {
						ServiceName   string `json:"service_name"`
						ServiceModule string `json:"service_module"`
						KeyVals       []struct {
							Key   string `json:"key"`
							Value string `json:"value"`
							Type  string `json:"type"`
						} `json:"key_vals"`
						Repos []struct {
							Source        string `json:"source"`
							RepoOwner     string `json:"repo_owner"`
							RepoNamespace string `json:"repo_namespace"`
							RepoName      string `json:"repo_name"`
							RemoteName    string `json:"remote_name"`
							Branch        string `json:"branch"`
							CodehostID    int    `json:"codehost_id"`
						} `json:"repos"`
					} `json:"default_service_and_builds"`
					ServiceAndBuilds []struct {
						ServiceName   string `json:"service_name"`
						ServiceModule string `json:"service_module"`
						BuildName     string `json:"build_name"`
						KeyVals       []struct {
							Key   string `json:"key"`
							Value string `json:"value"`
							Type  string `json:"type"`
						} `json:"key_vals"`
						Repos []struct {
							Source        string `json:"source"`
							RepoOwner     string `json:"repo_owner"`
							RepoNamespace string `json:"repo_namespace"`
							RepoName      string `json:"repo_name"`
							RemoteName    string `json:"remote_name"`
							Branch        string `json:"branch"`
							CodehostID    int    `json:"codehost_id"`
						} `json:"repos"`
					} `json:"service_and_builds"`
				} `json:"spec"`
			} `json:"jobs"`
		}{
			{
				Name: "构建",
				Jobs: []struct {
					Name string `json:"name"`
					Type string `json:"type"`
					Spec struct {
						DefaultServiceAndBuilds []struct {
							ServiceName   string `json:"service_name"`
							ServiceModule string `json:"service_module"`
							KeyVals       []struct {
								Key   string `json:"key"`
								Value string `json:"value"`
								Type  string `json:"type"`
							} `json:"key_vals"`
							Repos []struct {
								Source        string `json:"source"`
								RepoOwner     string `json:"repo_owner"`
								RepoNamespace string `json:"repo_namespace"`
								RepoName      string `json:"repo_name"`
								RemoteName    string `json:"remote_name"`
								Branch        string `json:"branch"`
								CodehostID    int    `json:"codehost_id"`
							} `json:"repos"`
						} `json:"default_service_and_builds"`
						ServiceAndBuilds []struct {
							ServiceName   string `json:"service_name"`
							ServiceModule string `json:"service_module"`
							BuildName     string `json:"build_name"`
							KeyVals       []struct {
								Key   string `json:"key"`
								Value string `json:"value"`
								Type  string `json:"type"`
							} `json:"key_vals"`
							Repos []struct {
								Source        string `json:"source"`
								RepoOwner     string `json:"repo_owner"`
								RepoNamespace string `json:"repo_namespace"`
								RepoName      string `json:"repo_name"`
								RemoteName    string `json:"remote_name"`
								Branch        string `json:"branch"`
								CodehostID    int    `json:"codehost_id"`
							} `json:"repos"`
						} `json:"service_and_builds"`
					} `json:"spec"`
				}{
					{
						Name: "构建发布",
						Type: "zadig-build",
						Spec: struct {
							DefaultServiceAndBuilds []struct {
								ServiceName   string `json:"service_name"`
								ServiceModule string `json:"service_module"`
								KeyVals       []struct {
									Key   string `json:"key"`
									Value string `json:"value"`
									Type  string `json:"type"`
								} `json:"key_vals"`
								Repos []struct {
									Source        string `json:"source"`
									RepoOwner     string `json:"repo_owner"`
									RepoNamespace string `json:"repo_namespace"`
									RepoName      string `json:"repo_name"`
									RemoteName    string `json:"remote_name"`
									Branch        string `json:"branch"`
									CodehostID    int    `json:"codehost_id"`
								} `json:"repos"`
							} `json:"default_service_and_builds"`
							ServiceAndBuilds []struct {
								ServiceName   string `json:"service_name"`
								ServiceModule string `json:"service_module"`
								BuildName     string `json:"build_name"`
								KeyVals       []struct {
									Key   string `json:"key"`
									Value string `json:"value"`
									Type  string `json:"type"`
								} `json:"key_vals"`
								Repos []struct {
									Source        string `json:"source"`
									RepoOwner     string `json:"repo_owner"`
									RepoNamespace string `json:"repo_namespace"`
									RepoName      string `json:"repo_name"`
									RemoteName    string `json:"remote_name"`
									Branch        string `json:"branch"`
									CodehostID    int    `json:"codehost_id"`
								} `json:"repos"`
							} `json:"service_and_builds"`
						}{
							DefaultServiceAndBuilds: []struct {
								ServiceName   string `json:"service_name"`
								ServiceModule string `json:"service_module"`
								KeyVals       []struct {
									Key   string `json:"key"`
									Value string `json:"value"`
									Type  string `json:"type"`
								} `json:"key_vals"`
								Repos []struct {
									Source        string `json:"source"`
									RepoOwner     string `json:"repo_owner"`
									RepoNamespace string `json:"repo_namespace"`
									RepoName      string `json:"repo_name"`
									RemoteName    string `json:"remote_name"`
									Branch        string `json:"branch"`
									CodehostID    int    `json:"codehost_id"`
								} `json:"repos"`
							}{
								{
									ServiceName:   sn,
									ServiceModule: sm,
									KeyVals: []struct {
										Key   string `json:"key"`
										Value string `json:"value"`
										Type  string `json:"type"`
									}{{
										Key:   "ENV_NAME",
										Value: "{{.workflow.params.环境}}",
										Type:  "string",
									}},
									Repos: []struct {
										Source        string `json:"source"`
										RepoOwner     string `json:"repo_owner"`
										RepoNamespace string `json:"repo_namespace"`
										RepoName      string `json:"repo_name"`
										RemoteName    string `json:"remote_name"`
										Branch        string `json:"branch"`
										CodehostID    int    `json:"codehost_id"`
									}{
										{
											Source:        "github",
											RepoOwner:     "storehubnet",
											RepoNamespace: "storehubnet",
											RepoName:      rn,
											RemoteName:    "origin",
											Branch:        req.BranchName,
											CodehostID:    6,
										},
									},
								},
							},
							ServiceAndBuilds: []struct {
								ServiceName   string `json:"service_name"`
								ServiceModule string `json:"service_module"`
								BuildName     string `json:"build_name"`
								KeyVals       []struct {
									Key   string `json:"key"`
									Value string `json:"value"`
									Type  string `json:"type"`
								} `json:"key_vals"`
								Repos []struct {
									Source        string `json:"source"`
									RepoOwner     string `json:"repo_owner"`
									RepoNamespace string `json:"repo_namespace"`
									RepoName      string `json:"repo_name"`
									RemoteName    string `json:"remote_name"`
									Branch        string `json:"branch"`
									CodehostID    int    `json:"codehost_id"`
								} `json:"repos"`
							}{
								{
									ServiceName:   sn,
									ServiceModule: sm,
									BuildName:     fmt.Sprintf("fat-base-envrionment-build-%s-1", bn),
									KeyVals: []struct {
										Key   string `json:"key"`
										Value string `json:"value"`
										Type  string `json:"type"`
									}{{
										Key:   "ENV_NAME",
										Value: "{{.workflow.params.环境}}",
										Type:  "string",
									}},
									Repos: []struct {
										Source        string `json:"source"`
										RepoOwner     string `json:"repo_owner"`
										RepoNamespace string `json:"repo_namespace"`
										RepoName      string `json:"repo_name"`
										RemoteName    string `json:"remote_name"`
										Branch        string `json:"branch"`
										CodehostID    int    `json:"codehost_id"`
									}{
										{
											Source:        "github",
											RepoOwner:     "storehubnet",
											RepoNamespace: "storehubnet",
											RepoName:      rn,
											RemoteName:    "origin",
											Branch:        req.BranchName,
											CodehostID:    6,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	ret := struct {
		ProjectName  string `json:"project_name"`
		WorkflowName string `json:"workflow_name"`
		TaskID       int    `json:"task_id"`
	}{}
	resp, err := z.client.R().
		SetResult(&ret).
		SetBody(servicesReq).
		Post("/api/aslan/workflow/v4/workflowtask?projectName=fat-base-envrionment")
	if err != nil {
		z.log.Warn("deploy srv err.", zap.Error(err))
		return 0, fmt.Errorf("deploy srv err: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		z.log.Warn("resp status code not ok",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response body", resp.String()),
		)
		return 0, fmt.Errorf("deploy srv status code not 200: %w", err)
	}

	return ret.TaskID, nil
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
		env, err := z.CreateSubEnv()
		if err != nil {
			z.log.Warn("create namespace failed.", zap.String("namespace", env))
			return err
		}
		taskId, err := z.DeployService(model.DeployServiceReq{
			SubEnv:      env,
			ServiceName: "backoffice-v1-web",
			BranchName:  "feature/INF-666",
			GithubActor: "jeremy2566",
		})
		if err != nil {
			z.log.Warn("deploy service failed.", zap.String("service", "backoffice-v1-web-app"))
		} else {
			z.log.Info("deploy service success.", zap.String("service_name", "backoffice-v1-web-app"), zap.Int("taskId", taskId))
		}
		// 加入 redis 缓存
		if err := z.rdb.SaveNamespace(env, model.DaoNamespace{
			SubEnv:      env,
			UpdateBy:    req.GithubActor,
			Branch:      branchName,
			ServiceName: []string{"backoffice-v1-web-app"},
		}); err != nil {
			z.log.Warn("save namespace failed.", zap.String("namespace", env))
		} else {
			z.log.Info("deploy service success.", zap.String("service_name", req.ServiceName))
		}
		err = z.AddService(model.AddServiceReq{
			SubEnv:      env,
			ServiceName: req.ServiceName,
		})
		// 添加新服务到 ns
		if err != nil {
			z.log.Warn("add service failed.", zap.String("service", req.ServiceName))
			return err
		}
		taskId, err = z.DeployService(model.DeployServiceReq{
			SubEnv:      env,
			ServiceName: req.ServiceName,
			BranchName:  branchName,
			GithubActor: "jeremy2566",
		})
		if err != nil {
			z.log.Warn("deploy service failed.", zap.String("service_name", req.ServiceName), zap.Int("taskId", taskId))
		} else {
			z.log.Info("deploy service success.", zap.String("service_name", req.ServiceName), zap.Int("taskId", taskId))
		}
		z.rdb.UpdateServiceByKey(env, req.ServiceName)
	} else {
		// 添加新服务到 ns
		err := z.AddService(model.AddServiceReq{
			SubEnv:      namespace.SubEnv,
			ServiceName: req.ServiceName,
		})
		if err != nil {
			z.log.Warn("add service failed.", zap.String("service", req.ServiceName))
			return err
		}
		taskId, err := z.DeployService(model.DeployServiceReq{
			SubEnv:      namespace.SubEnv,
			ServiceName: req.ServiceName,
			BranchName:  branchName,
			GithubActor: "jeremy2566",
		})
		if err != nil {
			z.log.Warn("deploy service failed.", zap.String("service_name", req.ServiceName), zap.Int("taskId", taskId))
		} else {
			z.log.Info("deploy service success.", zap.String("service_name", req.ServiceName), zap.Int("taskId", taskId))
		}
		z.rdb.UpdateServiceByKey(namespace.SubEnv, req.ServiceName)
	}
	return nil
}

func (z *zadigImpl) CreateSubEnv() (string, error) {
	namespace := z.selectedNamespace()
	se := model.ShareEnvReq{
		Enable:  true,
		IsBase:  false,
		BaseEnv: "test33",
	}
	var cvs []model.ChartInfoReq
	charts := z.GetServiceCharts()
	for _, service := range []string{"redis-backoffice", "redis-general", "backoffice-v1-web-app", "bo-v1-assets"} {
		value, exist := charts[service]
		if !exist {
			z.log.Warn("service not found", zap.String("service", service))
			continue
		}

		cvs = append(cvs, model.ChartInfoReq{
			EnvName:        namespace,
			ServiceName:    service,
			ChartVersion:   value,
			DeployStrategy: "deploy",
		})
	}
	req := model.CreateSubEnvReq{
		{
			EnvName:     namespace,
			ClusterID:   "64e48b5b8fc410571753cc6c",
			RegistryID:  "64e485c78fc410571753cc67",
			ChartValues: cvs,
			Namespace:   namespace,
			IsExisted:   false,
			ShareEnv:    se,
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
	// 等待环境创建完成
	if err := z.waitForEnvReady(namespace, 5*time.Minute); err != nil {
		z.log.Error("environment creation failed", zap.String("namespace", namespace), zap.Error(err))
		return namespace, fmt.Errorf("environment creation failed: %w", err)
	}

	return namespace, nil
}

// waitForEnvReady 等待环境状态变为success
func (z *zadigImpl) waitForEnvReady(envName string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second) // 每5秒检查一次
	defer ticker.Stop()

	z.log.Info("waiting for environment to be ready", zap.String("env", envName))

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for environment %s to be ready", envName)
		case <-ticker.C:
			detail, err := z.GetTestEnvDetail(envName, "fat-base-envrionment")
			if err != nil {
				z.log.Warn("failed to get env detail", zap.String("env", envName), zap.Error(err))
				continue
			}

			z.log.Info("checking environment status", zap.String("env", envName), zap.String("status", detail.Status))

			if detail.Status == "success" {
				z.log.Info("environment is ready", zap.String("env", envName))
				return nil
			}

			if detail.Status == "error" {
				return fmt.Errorf("environment %s creation failed with status: %s", envName, detail.Status)
			}
		}
	}
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
		ChartInfos []model.RespChartInfo `json:"chart_infos"`
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
