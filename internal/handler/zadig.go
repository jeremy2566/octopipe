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
	}
	if resp.StatusCode() != http.StatusOK {
		h.log.Warn("resp status code not ok",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response body", resp.String()),
		)
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

func (h Handler) AddServices(namespace string, servicesMap map[string]string) error {
	h.log.Info("add services", zap.String("namespace", namespace), zap.Any("servicesMap", servicesMap))

	var services []AddServiceList
	for sn := range servicesMap {
		add := AddServiceList{ServiceName: sn}
		services = append(services, add)
	}
	req := AddService{
		EnvKey:   namespace,
		Services: services,
	}
	apiResponse := struct {
		Message string `json:"message"`
	}{}
	resp, err := h.client.R().
		SetResult(&apiResponse).
		SetContentType("application/json").
		SetAuthToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiamVyZW15MjU2NiIsImVtYWlsIjoiamVyZW15LnpoYW5nQHN0b3JlaHViLmNvbSIsInVpZCI6Ijk3ODgyYzVmLWEyNjYtMTFlZi1hYTlmLTAyMDU4ZWVlYTIzNSIsInByZWZlcnJlZF91c2VybmFtZSI6ImplcmVteTI1NjYiLCJmZWRlcmF0ZWRfY2xhaW1zIjp7ImNvbm5lY3Rvcl9pZCI6ImdpdGh1YiIsInVzZXJfaWQiOiJqZXJlbXkyNTY2In0sImF1ZCI6InphZGlnIiwiZXhwIjo0ODg1MTc0NzMwfQ.pZ_jVTj20h_R9Z84O3_QJL2OcUxzJn04gNkDIATRsf4").
		SetBody(req).
		Post("https://zadigx.shub.us/openapi/environments/service/yaml?projectKey=fat-base-envrionment")

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

func (h Handler) DeployServices(namespace string, servicesMap map[string]string) error {
	h.log.Info("deploy services", zap.String("namespace", namespace), zap.Any("servicesMap", servicesMap))
	// TODO 部署服务
	return nil
}
