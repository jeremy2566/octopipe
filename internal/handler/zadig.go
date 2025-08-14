package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeremy2566/octopipe/internal/cache"
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
	sn, sm, rn, bn := trans(serviceName)
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
									ServiceName:   sn,
									ServiceModule: sm,
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
									ServiceName:   sn,
									ServiceModule: sm,
									BuildName:     fmt.Sprintf("fat-base-envrionment-build-%s-1", bn),
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

	marshal, _ := json.Marshal(req)
	fmt.Println(string(marshal))

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

// trans 服务名转 service name, service module, repo name 和 build name
func trans(serviceName string) (string, string, string, string) {
	switch serviceName {
	case "backoffice-v1-web":
		return "backoffice-v1-web-app", "backoffice-v1-web", "backoffice-v1-web", "backoffice-v1-web"
	default:
		return serviceName, serviceName, serviceName, serviceName
	}
}

type callback struct {
	ObjectKind string `json:"object_kind"`
	Event      string `json:"event"`
	Workflow   struct {
		TaskID              int    `json:"task_id"`
		ProjectName         string `json:"project_name"`
		ProjectDisplayName  string `json:"project_display_name"`
		WorkflowName        string `json:"workflow_name"`
		WorkflowDisplayName string `json:"workflow_display_name"`
		Status              string `json:"status"`
		Remark              string `json:"remark"`
		DetailURL           string `json:"detail_url"`
		Error               string `json:"error"`
		CreateTime          int    `json:"create_time"`
		StartTime           int    `json:"start_time"`
		EndTime             int    `json:"end_time"`
		Stages              []struct {
			Name      string `json:"name"`
			Status    string `json:"status"`
			StartTime int    `json:"start_time"`
			EndTime   int    `json:"end_time"`
			Jobs      []struct {
				Name        string `json:"name"`
				DisplayName string `json:"display_name"`
				Type        string `json:"type"`
				Status      string `json:"status"`
				StartTime   int    `json:"start_time"`
				EndTime     int    `json:"end_time"`
				Error       string `json:"error"`
				Spec        struct {
					Repositories []struct {
						Source        string      `json:"source"`
						RepoOwner     string      `json:"repo_owner"`
						RepoNamespace string      `json:"repo_namespace"`
						RepoName      string      `json:"repo_name"`
						Branch        string      `json:"branch"`
						Prs           interface{} `json:"prs"`
						Tag           string      `json:"tag"`
						CommitID      string      `json:"commit_id"`
						CommitURL     string      `json:"commit_url"`
						CommitMessage string      `json:"commit_message"`
					} `json:"repositories"`
					Image string `json:"image"`
				} `json:"spec"`
			} `json:"jobs"`
			Error string `json:"error"`
		} `json:"stages"`
		TaskCreator      string `json:"task_creator"`
		TaskCreatorID    string `json:"task_creator_id"`
		TaskCreatorPhone string `json:"task_creator_phone"`
		TaskCreatorEmail string `json:"task_creator_email"`
		TaskType         string `json:"task_type"`
	} `json:"workflow"`
}

func (h Handler) Webhook(c *gin.Context) {
	var cb callback
	if err := c.ShouldBindJSON(&cb); err != nil {
		h.log.Warn("bind json failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Info("received request for zadig webhook.", zap.Any("params", cb))

	var success bool
	switch cb.Workflow.Status {
	case "passed":
		success = true
	default:
		success = false
	}
	branch := cb.Workflow.Stages[0].Jobs[0].Spec.Repositories[0].Branch

	var subEnv string
	if ns, exist := cache.GetInstance().Get(branch); exist {
		subEnv = ns
	} else {
		subEnv = "test-1"
	}
	totalTime := cb.Workflow.EndTime - cb.Workflow.CreateTime
	req := SenderReq{
		ProjectName:    cb.Workflow.ProjectName,
		WorkflowName:   cb.Workflow.WorkflowDisplayName,
		WorkflowNumber: cb.Workflow.TaskID,
		Duration:       fmt.Sprintf("%02d:%02d", totalTime/60, totalTime%60),
		SubEnv:         subEnv,
		Service:        cb.Workflow.Stages[0].Jobs[0].Spec.Repositories[0].RepoName,
		Branch:         branch,
		Success:        success,
		Email:          cb.Workflow.TaskCreatorEmail,
	}
	h.Sender(req)
}
