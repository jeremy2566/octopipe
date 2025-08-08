package handler

import (
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type env struct {
	EnvKey     string `json:"env_key"`
	ClusterID  string `json:"cluster_id"`
	Namespace  string `json:"namespace"`
	Production bool   `json:"production"`
	RegistryID string `json:"registry_id"`
	Status     string `json:"status"`
	UpdateBy   string `json:"update_by"`
	UpdateTime int64  `json:"update_time"`
}

func (h Handler) GetZadigNamespace(projectKey string) []env {
	var envs []env
	resp, err := h.client.R().
		SetResult(&envs).
		SetQueryParam("projectKey", projectKey).
		SetContentType("application/json").
		SetAuthToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiamVyZW15MjU2NiIsImVtYWlsIjoiamVyZW15LnpoYW5nQHN0b3JlaHViLmNvbSIsInVpZCI6Ijk3ODgyYzVmLWEyNjYtMTFlZi1hYTlmLTAyMDU4ZWVlYTIzNSIsInByZWZlcnJlZF91c2VybmFtZSI6ImplcmVteTI1NjYiLCJmZWRlcmF0ZWRfY2xhaW1zIjp7ImNvbm5lY3Rvcl9pZCI6ImdpdGh1YiIsInVzZXJfaWQiOiJqZXJlbXkyNTY2In0sImF1ZCI6InphZGlnIiwiZXhwIjo0ODg1MTc0NzMwfQ.pZ_jVTj20h_R9Z84O3_QJL2OcUxzJn04gNkDIATRsf4").
		Get("https://zadigx.shub.us/openapi/environments")
	if err != nil {
		h.log.Warn("failed to fetch environments", zap.Error(err))
		return nil
	}
	if resp.StatusCode() != http.StatusOK {
		h.log.Warn("resp status code not ok",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response body", resp.String()),
		)
		return nil
	}

	return envs
}

type CreateSubEnvReq []struct {
	EnvName     string        `json:"env_name"`
	ClusterID   string        `json:"cluster_id"`
	RegistryID  string        `json:"registry_id"`
	ChartValues []ChartValues `json:"chartValues"`
	Namespace   string        `json:"namespace"`
	IsExisted   bool          `json:"is_existed"`
	ShareEnv    ShareEnv      `json:"share_env"`
}
type ChartValues struct {
	EnvName        string `json:"envName"`
	ServiceName    string `json:"serviceName"`
	ChartVersion   string `json:"chartVersion"`
	DeployStrategy string `json:"deploy_strategy"`
}
type ShareEnv struct {
	Enable  bool   `json:"enable"`
	IsBase  bool   `json:"isBase"`
	BaseEnv string `json:"base_env"`
}

func (h Handler) CreateSubEnv(namespace string, services ...string) error {
	if !strings.HasPrefix(namespace, "test") {
		err := fmt.Errorf("namespace prefix cannot be 'test', got: %s", namespace)
		h.log.Error("validation failed for namespace",
			zap.String("namespace", namespace),
			zap.Error(err),
		)
		return err
	}
	if len(services) <= 0 {
		err := fmt.Errorf("services is empty")
		h.log.Error("validation failed for services",
			zap.Error(err),
		)
		return err
	}
	var cvs []ChartValues
	charts := h.GetServiceCharts()
	for _, service := range services {
		value, exist := charts[service]
		if !exist {
			h.log.Warn("service not found", zap.String("service", service))
			continue
		}

		cvs = append(cvs, ChartValues{
			EnvName:        namespace,
			ServiceName:    service,
			ChartVersion:   value,
			DeployStrategy: "deploy",
		})
	}

	se := ShareEnv{
		Enable:  true,
		IsBase:  false,
		BaseEnv: "test33",
	}

	req := CreateSubEnvReq{
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

	resp, err := h.client.R().
		SetResult(&apiResponse).
		SetContentType("application/json").
		SetAuthToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiamVyZW15MjU2NiIsImVtYWlsIjoiamVyZW15LnpoYW5nQHN0b3JlaHViLmNvbSIsInVpZCI6Ijk3ODgyYzVmLWEyNjYtMTFlZi1hYTlmLTAyMDU4ZWVlYTIzNSIsInByZWZlcnJlZF91c2VybmFtZSI6ImplcmVteTI1NjYiLCJmZWRlcmF0ZWRfY2xhaW1zIjp7ImNvbm5lY3Rvcl9pZCI6ImdpdGh1YiIsInVzZXJfaWQiOiJqZXJlbXkyNTY2In0sImF1ZCI6InphZGlnIiwiZXhwIjo0ODg1MTc0NzMwfQ.pZ_jVTj20h_R9Z84O3_QJL2OcUxzJn04gNkDIATRsf4").
		SetBody(req).
		Post("https://zadigx.shub.us/api/aslan/environment/environments?type=helm&projectName=fat-base-envrionment")
	if err != nil {
		h.log.Warn("failed to fetch environments", zap.Error(err))
	}
	if resp.StatusCode() != http.StatusOK {
		h.log.Warn("resp status code not ok",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response body", resp.String()),
		)
	}

	return nil
}

type ChartInfo struct {
	ServiceName  string `json:"service_name"`
	ChartVersion string `json:"chart_version"`
}

func (h Handler) GetServiceCharts() map[string]string {
	apiResponse := struct {
		ChartInfos []ChartInfo `json:"chart_infos"`
	}{}
	resp, err := h.client.R().
		SetResult(&apiResponse).
		SetAuthToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiamVyZW15MjU2NiIsImVtYWlsIjoiamVyZW15LnpoYW5nQHN0b3JlaHViLmNvbSIsInVpZCI6Ijk3ODgyYzVmLWEyNjYtMTFlZi1hYTlmLTAyMDU4ZWVlYTIzNSIsInByZWZlcnJlZF91c2VybmFtZSI6ImplcmVteTI1NjYiLCJmZWRlcmF0ZWRfY2xhaW1zIjp7ImNvbm5lY3Rvcl9pZCI6ImdpdGh1YiIsInVzZXJfaWQiOiJqZXJlbXkyNTY2In0sImF1ZCI6InphZGlnIiwiZXhwIjo0ODg1MTc0NzMwfQ.pZ_jVTj20h_R9Z84O3_QJL2OcUxzJn04gNkDIATRsf4").
		Get("https://zadigx.shub.us/api/aslan/environment/init_info/fat-base-envrionment?envType=share&isBaseEnv=false&baseEnv=test33&projectName=fat-base-envrionment")
	if err != nil {
		h.log.Warn("failed to fetch environments", zap.Error(err))
	}
	if resp.StatusCode() != http.StatusOK {
		h.log.Warn("resp status code not ok",
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

type AddService struct {
	EnvKey   string           `json:"env_key"`
	Services []AddServiceList `json:"service_list"`
}

type AddServiceList struct {
	ServiceName string `json:"service_name"`
}

type UtilsFun struct {
	ReplacePolicy string                `json:"replacePolicy"`
	EnvNames      []string              `json:"envNames"`
	ChartValues   []UtilsFunChartValues `json:"chartValues"`
}

type UtilsFunChartValues struct {
	EnvName         string `json:"envName"`
	ServiceName     string `json:"serviceName"`
	ReleaseName     string `json:"releaseName"`
	ChartVersion    string `json:"chartVersion"`
	Deploy_strategy string `json:"deploy_strategy"`
}

func (h Handler) AddServices(namespace, serviceName, branchName string) error {
	h.log.Info("add services", zap.String("namespace", namespace), zap.String("serviceName", serviceName), zap.String("branchName", branchName))
	charts := h.GetServiceCharts()
	chartVersion, exist := charts[serviceName]
	if !exist {
		h.log.Warn("service not found", zap.String("service", serviceName))
		return fmt.Errorf("service not found")
	}
	ufcv := []UtilsFunChartValues{
		{
			EnvName:         namespace,
			ServiceName:     serviceName,
			ReleaseName:     serviceName,
			ChartVersion:    chartVersion,
			Deploy_strategy: "deploy",
		},
	}
	uf := UtilsFun{
		ReplacePolicy: "notUseEnvImage",
		EnvNames:      []string{namespace},
		ChartValues:   ufcv,
	}

	var apiResponse []struct {
		EnvName    string `json:"env_name"`
		Status     string `json:"status"`
		ErrMessage string `json:"err_message"`
	}
	resp, err := h.client.R().
		SetResult(&apiResponse).
		SetContentType("application/json").
		SetAuthToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiamVyZW15MjU2NiIsImVtYWlsIjoiamVyZW15LnpoYW5nQHN0b3JlaHViLmNvbSIsInVpZCI6Ijk3ODgyYzVmLWEyNjYtMTFlZi1hYTlmLTAyMDU4ZWVlYTIzNSIsInByZWZlcnJlZF91c2VybmFtZSI6ImplcmVteTI1NjYiLCJmZWRlcmF0ZWRfY2xhaW1zIjp7ImNvbm5lY3Rvcl9pZCI6ImdpdGh1YiIsInVzZXJfaWQiOiJqZXJlbXkyNTY2In0sImF1ZCI6InphZGlnIiwiZXhwIjo0ODg1MTc0NzMwfQ.pZ_jVTj20h_R9Z84O3_QJL2OcUxzJn04gNkDIATRsf4").
		SetBody(uf).
		Put("https://zadigx.shub.us/api/aslan/environment/environments?type=helm&projectName=fat-base-envrionment")

	if err != nil {
		h.log.Warn("failed to fetch environments", zap.Error(err))
	}
	if resp.StatusCode() != http.StatusOK {
		h.log.Warn("resp status code not ok",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response body", resp.String()),
		)
		//return fmt.Errorf(resp.String())
	}

	return nil
}

type DeployServicesReq struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Params      []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"params"`
	Stages []struct {
		Name string `json:"name"`
		Jobs []struct {
			Name string `json:"name"`
			Type string `json:"type"`
			Spec struct {
				DefaultServiceAndBuilds []struct {
					ServiceName   string `json:"service_name"`
					ServiceModule string `json:"service_module"`
					Repos         []struct {
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
					Repos         []struct {
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
	} `json:"stages"`
	Project string `json:"project"`
}

func (h Handler) DeployService(namespace, serviceName, branchName string) (int, error) {
	h.log.Info("deploy services",
		zap.String("namespace", namespace),
		zap.String("serviceName", serviceName),
		zap.String("branchName", branchName),
	)
	req := DeployServicesReq{
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
				Value: namespace,
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
						Repos         []struct {
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
						Repos         []struct {
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
							Repos         []struct {
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
							Repos         []struct {
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
								Repos         []struct {
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
								Repos         []struct {
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
								Repos         []struct {
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
									ServiceName:   serviceName,
									ServiceModule: serviceName,
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
											RepoName:      serviceName,
											RemoteName:    "origin",
											Branch:        branchName,
											CodehostID:    6,
										},
									},
								},
							},
							ServiceAndBuilds: []struct {
								ServiceName   string `json:"service_name"`
								ServiceModule string `json:"service_module"`
								BuildName     string `json:"build_name"`
								Repos         []struct {
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
									ServiceName:   serviceName,
									ServiceModule: serviceName,
									BuildName:     fmt.Sprintf("fat-base-envrionment-build-%s-1", serviceName),
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
											RepoName:      serviceName,
											RemoteName:    "origin",
											Branch:        branchName,
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
	resp, err := h.client.R().
		SetResult(&ret).
		SetContentType("application/json").
		SetAuthToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiamVyZW15MjU2NiIsImVtYWlsIjoiamVyZW15LnpoYW5nQHN0b3JlaHViLmNvbSIsInVpZCI6Ijk3ODgyYzVmLWEyNjYtMTFlZi1hYTlmLTAyMDU4ZWVlYTIzNSIsInByZWZlcnJlZF91c2VybmFtZSI6ImplcmVteTI1NjYiLCJmZWRlcmF0ZWRfY2xhaW1zIjp7ImNvbm5lY3Rvcl9pZCI6ImdpdGh1YiIsInVzZXJfaWQiOiJqZXJlbXkyNTY2In0sImF1ZCI6InphZGlnIiwiZXhwIjo0ODg1MTc0NzMwfQ.pZ_jVTj20h_R9Z84O3_QJL2OcUxzJn04gNkDIATRsf4").
		SetBody(req).
		Post("https://zadigx.shub.us/api/aslan/workflow/v4/workflowtask?projectName=fat-base-envrionment")
	if err != nil {
		h.log.Warn("deploy srv err.", zap.Error(err))
		return -1, fmt.Errorf("deploy srv err: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		h.log.Warn("resp status code not ok",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response body", resp.String()),
		)
		return -1, fmt.Errorf("deploy srv status code not 200: %w", err)
	}
	return ret.TaskID, nil
}
