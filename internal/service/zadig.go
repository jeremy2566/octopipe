package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
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
	Webhook(cb model.Callback) error
	GetTaskDetail(workflowKey string, taskId int) (model.RespTaskDetail, error)
}

type zadigImpl struct {
	log    *zap.Logger
	client *resty.Client
	rdb    dao.Rdb
	lark   Lark
}

func (z *zadigImpl) GetTaskDetail(workflowKey string, taskId int) (model.RespTaskDetail, error) {
	ret := model.RespTaskDetail{}
	resp, err := z.client.R().
		SetQueryParams(map[string]string{
			"workflowKey": workflowKey,
			"taskId":      strconv.Itoa(taskId),
		}).
		SetResult(&ret).
		Get("openapi/workflows/custom/task")

	if err != nil {
		z.log.Warn("failed to fetch workflow detail", zap.Error(err))
		return ret, err
	}
	if resp.StatusCode() != http.StatusOK {
		z.log.Warn("resp status code not ok",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response body", resp.String()),
		)
	}
	return ret, nil
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

// trans 服务名转 service_name, service_module, repo_name 和 build_name
func (z *zadigImpl) trans(serviceName string) (string, string, string, string) {
	switch serviceName {
	case "accounting-cleartax-connector":
		return "account-cleartax-connector", "account-cleartax-connector", "accounting-cleartax-connector", "account-cleartax-connector"
	case "app_push-infra-svc":
		return "apppush-infra", "apppush-infra", "app_push-infra-svc", "fat-base-envrionment-build-apppush-infra-1"
	case "auth-api":
		return "auth-api", "auth-api", "auth-api", "fat-base-envrionment-build-auth-api-1"
	case "backoffice-v2-bff":
		return "backoffice-v2-bff", "backoffice-v2-bff", "backoffice-v2-bff", "fat-base-envrionment-build-backoffice-v2-bff-1"
	case "beep-v1-web":
		return "beep-v1-web", "beep-v1-web", "beep-v1-web", "fat-base-envrionment-build-beep-v1-web-1"
	case "campaign-svc":
		return "campaign-svc", "campaign-svc", "campaign-svc", "fat-base-envrionment-build-campaign-svc-1"
	case "core-api":
		return "core-api", "core-api", "core-api", "fat-base-envrionment-build-core-api-1"
	case "crm-api":
		return "crm-api", "crm-api", "crm-api", "fat-base-envrionment-build-crm-api-1"
	case "customer-svc":
		return "customer-svc", "customer-svc", "customers-svc", "fat-base-envrionment-build-customer-svc-1"
	case "e-invoice_adapter-svc":
		return "e-invoice-adapter-svc", "e-invoice-adapter-svc", "e-invoice_adapter-svc", "e-invoice-test33"
	case "ecommerce-v1-api":
		return "ecommerce-v1-api", "ecommerce-v1-api", "ecommerce-v1-api", "fat-base-envrionment-build-ecommerce-v1-api-1"
	case "ecommerce-v1-web":
		return "ecommerce-v1-web", "ecommerce-v1-web", "ecommerce-v1-web", "fat-base-envrionment-build-ecommerce-v1-web-1"
	case "employee-domain-svc":
		return "employee-domain-svc", "employee-domain-svc", "employee-domain-svc", "fat-base-envrionment-build-employee-domain-svc-1"
	case "ist-v1-web":
		return "ist-v1-web", "ist-v1-web", "ist-v1-web", "fat-base-envrionment-build-ist-v1-web-1"
	case "logistics-app-svc":
		return "logistics-app-svc", "logistics-app-svc", "logistics-app-svc", "fat-base-envrionment-build-logistics-app-svc-1"
	case "merchant-domain-svc":
		return "merchant-domain-svc", "merchant-domain-svc", "merchant-domain-svc", "fat-base-envrionment-build-merchant-domain-svc-1"
	case "messenger_apps-infra-svc":
		return "messenger-apps-infra-svc", "messenger-apps-infra-svc", "messenger_apps-infra-svc", "fat-base-envrionment-build-messenger-apps-infra-svc-1"
	case "online_store-domain-svc":
		return "online-store-domain-svc", "online-store-domain-svc", "online_store-domain-svc", "fat-base-envrionment-build-online-store-domain-svc-1"
	case "ost-v1-web":
		return "ost-v1-web", "ost-v1-web", "ost-v1-web", "fat-base-envrionment-build-ost-v1-web-1"
	case "otp-api":
		return "otp-api", "otp-api", "otp-api", "fat-base-envrionment-build-otp-api-1"
	case "payment-api":
		return "payment-api", "payment-api", "payment-api", "fat-base-envrionment-build-payment-api-1"
	case "payout-api":
		return "payout-api", "payout-api", "payout-api", "fat-base-envrionment-build-payout-api-1"
	case "pos-v3-bff":
		return "pos-v3-bff", "pos-v3-bff", "pos-v3-bff", "fat-base-envrionment-build-pos-v3-bff-1"
	case "product-domain-svc":
		return "product-domain-svc", "product-domain-svc", "product-domain-svc", "fat-base-envrionment-build-product-domain-svc-1"
	case "promotion-api":
		return "promotion-api", "promotion-api", "promotion-api", "fat-base-envrionment-build-promotion-api-1"
	case "promotion-svc":
		return "promotion-svc", "promotion-svc", "promotion-svc", "fat-base-envrionment-build-promotion-svc-1"
	case "realtime_event_broker-svc":
		return "realtime-event-broker-svc", "realtime-event-broker-svc", "realtime_event_broker-svc", "realtime-event-broker-svc"
	case "report-api":
		return "report-api", "report-api", "report-api", "fat-base-envrionment-build-report-api-1"
	case "shmanager-v1-bff":
		return "shmanager-v1-bff", "shmanager-v1-bff", "shmanager-v1-bff", "fat-base-envrionment-build-shmanager-v1-bff-1"
	case "sms-api":
		return "sms-api", "sms-api", "sms-api", "fat-base-envrionment-build-sms-api-1"
	case "subscription_adapter-infra-svc":
		return "sub-adapter-infra-svc", "sub-adapter-infra-svc", "subscription_adapter-infra-svc", "fat-base-envrionment-build-sub-adapter-infra-svc-1"
	case "subscription-app-svc":
		return "subscription-app-svc", "subscription-app-svc", "subscription-app-svc", "fat-base-envrionment-build-subscription-app-svc-1"
	case "3pl_adapter-infra-svc":
		return "th3pl-adapter-infra-svc", "th3pl-adapter-infra-svc", "3pl_adapter-infra-svc", "fat-base-envrionment-build-th3pl-adapter-infra-svc-1"
	case "3rd_party_food_delivery_adapter-svc":
		return "thpfd-adapter", "thpfd-adapter", "3rd_party_food_delivery_adapter-svc", "fat-base-envrionment-build-thpfd-adapter-1"
	case "user_preference-api":
		return "user-preference-api", "user-preference-api", "user_preference-api", "user-preference-api-test33"
	case "3p_webhook_adapter-infra-svc":
		return "webhook-adapter-infra-svc", "webhook-adapter-infra-svc", "3p_webhook_adapter-infra-svc", "fat-base-envrionment-build-webhook-adapter-infra-svc-1"
	case "webhook-svc":
		return "webhook-svc", "webhook-svc", "webhook-svc", "fat-base-envrionment-build-webhook-svc-1"
	default:
		return serviceName, serviceName, serviceName, serviceName
	}
}

// transCEC core-event-consumer 服务名转 service_name, service_module, repo_name 和 build_name
func (z *zadigImpl) transCEC(branchName string) (string, string, string, string) {
	switch {
	case strings.Contains(branchName, "beep"):
		return "core-event-consumer-beep", "core-event-consumer-beep", "core-event-consumer", "fat-base-envrionment-build-core-event-consumer-beep-1"
	case strings.Contains(branchName, "beepapp"):
		return "core-event-consumer-beepapp", "core-event-consumer-beepapp", "core-event-consumer", "fat-base-envrionment-build-core-event-consumer-beepapp-1"
	case strings.Contains(branchName, "bo"):
		return "core-event-consumer-bo", "core-event-consumer-bo", "core-event-consumer", "fat-base-envrionment-build-core-event-consumer-bo-1"
	case strings.Contains(branchName, "msg"):
		return "core-event-consumer-msg", "core-event-consumer-msg", "core-event-consumer", "core-event-consumer-msg-test33"
	case strings.Contains(branchName, "sms"):
		return "core-event-consumer-sms", "core-event-consumer-sms", "core-event-consumer", "core-event-consumer-sms-test33"
	case strings.Contains(branchName, "otpsms"):
		return "core-event-consumer-otpsms", "core-event-consumer-otpsms", "core-event-consumer", "core-event-consumer-otpsms-test33"
	case strings.Contains(branchName, "payment"):
		return "core-event-consumer-payment", "core-event-consumer-payment", "core-event-consumer", "fat-base-envrionment-build-core-event-consumer-payment-1"
	case strings.Contains(branchName, "zendesk"):
		return "core-event-consumer-zendesk", "core-event-consumer-zendesk", "core-event-consumer", "core-event-consumer-zendesk-test33"
	default:
		z.log.Warn("unknown branch name", zap.String("branchName", branchName))
		return "", "", "", ""
	}
}

// transOPS online_purchase-svc 服务名转 service_name, service_module, repo_name 和 build_name
func (z *zadigImpl) transOPS(branchName string) (string, string, string, string) {
	switch {
	case strings.Contains(branchName, "cronjob"):
		return "online-purchase-svc-cronjob", "online-purchase-svc-cronjob", "online_purchase-svc", "online-purchase-svc-cronjob-test33"
	default:
		return "online-purchase-svc", "online-purchase-svc", "online_purchase-svc", "fat-base-envrionment-build-online-purchase-svc-1"
	}
}

// transIS inventory-svc 服务名转 service_name, service_module, repo_name 和 build_name
func (z *zadigImpl) transIS(branchName string) (string, string, string, string) {
	switch {
	case strings.Contains(branchName, "stockjob"):
		return "inventory-stockjob", "inventory-domain-svc", "inventory-svc", "inventory-stockjob-test33"
	default:
		return "inventory-domain-svc", "inventory-domain-svc", "inventory-svc", "fat-base-envrionment-build-inventory-domain-svc-1"
	}
}

// transEVC ecommerce-v1-consumer 服务名转 service_name, service_module, repo_name 和 build_name
func (z *zadigImpl) transEVC(branchName string) (string, string, string, string) {
	switch {
	case strings.Contains(branchName, "checkout"):
		return "ec-v1-consumer-checkout", "ec-v1-consumer-checkout", "ecommerce-v1-consumer", "fat-base-envrionment-build-ec-v1-consumer-checkout-1"
	default:
		return "ec-v1-consumer", "ec-v1-consumer", "ecommerce-v1-consumer", "fat-base-envrionment-build-ec-v1-consumer-1"
	}
}

// transBoV1Web 服务名转 service_name, service_module, repo_name 和 build_name
func (z *zadigImpl) transBoV1Web(branchName string) (string, string, string, string) {
	switch {
	case strings.Contains(branchName, "addonsjob"):
		return "backoffice-v1-web-addonsjob", "backoffice-v1-web", "backoffice-v1-web", "backoffice-v1-web-addonsjob"
	case strings.Contains(branchName, "api"):
		return "backoffice-v1-web-api", "backoffice-v1-web", "backoffice-v1-web", "backoffice-v1-web-api-test33"
	case strings.Contains(branchName, "exportjob"):
		return "backoffice-v1-web-exportjob", "backoffice-v1-web-exportjob", "backoffice-v1-web", "fat-base-envrionment-build-backoffice-v1-web-exportjob"
	case strings.Contains(branchName, "importjob"):
		return "backoffice-v1-web-importjob", "backoffice-v1-web", "backoffice-v1-web", "backoffice-v1-web-importjob-test33"
	case strings.Contains(branchName, "migratejob"):
		return "backoffice-v1-web-migratejob", "backoffice-v1-web-migratejob", "backoffice-v1-web", "backoffice-v1-web-migratejob-test33"
	case strings.Contains(branchName, "qbojob"):
		return "backoffice-v1-web-qbojob", "backoffice-v1-web", "backoffice-v1-web", "fat-base-envrionment-build-backoffice-v1-web-1"
	case strings.Contains(branchName, "running"):
		return "backoffice-v1-web-runningjob", "backoffice-v1-web", "backoffice-v1-web", "backoffice-v1-web-runningjob-test33"
	case strings.Contains(branchName, "scheduledjob"):
		return "backoffice-v1-web-scheduledjob", "backoffice-v1-web", "backoffice-v1-web", "backoffice-v1-web-schedule-job"
	default:
		return "backoffice-v1-web-app", "backoffice-v1-web", "backoffice-v1-web", "fat-base-envrionment-build-backoffice-v1-web-1"
	}
}

func (z *zadigImpl) DeployService(req model.DeployServiceReq) (int, error) {
	z.log.Info("deploy service.", zap.Any("params", req))
	var sn, sm, rn, bn string
	if req.ServiceName == "backoffice-v1-web" {
		sn, sm, rn, bn = z.transBoV1Web(req.BranchName)
	} else if req.ServiceName == "core-event-consumer" {
		sn, sm, rn, bn = z.transCEC(req.BranchName)
	} else if req.ServiceName == "ecommerce-v1-consumer" {
		sn, sm, rn, bn = z.transEVC(req.BranchName)
	} else if req.ServiceName == "inventory-svc" {
		sn, sm, rn, bn = z.transIS(req.ServiceName)
	} else if req.ServiceName == "online_purchase-svc" {
		sn, sm, rn, bn = z.transOPS(req.ServiceName)
	} else {
		sn, sm, rn, bn = z.trans(req.ServiceName)
	}

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
									BuildName:     bn,
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
		lark: NewLark(
			log,
			resty.New().SetRetryCount(3).SetRetryWaitTime(1*time.Second).SetRetryMaxWaitTime(5*time.Second),
		),
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
		split := strings.Split(branchName, "/")
		num := split[len(split)-1]
		taskId, err := z.DeployService(model.DeployServiceReq{
			SubEnv:      env,
			ServiceName: "backoffice-v1-web",
			BranchName:  fmt.Sprintf("feature/app/%s", num),
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
	ret["accounting-cleartax-connector"] = "4"
	return ret
}

func (z *zadigImpl) Webhook(cb model.Callback) error {
	z.log.Info("received webhook callback",
		zap.String("workflow_name", cb.Workflow.WorkflowName),
		zap.String("status", cb.Workflow.Status),
		zap.Int("task_id", cb.Workflow.TaskID),
		zap.String("project", cb.Workflow.ProjectName))

	// 根据不同的 workflow name 路由到不同的处理器
	switch cb.Workflow.WorkflowName {
	case "domain-monitor":
		return z.domainMonitor(cb)
	case "test33":
		return z.deploySubEnv(cb)
	default:
		z.log.Warn("unknown workflow name, using default handler",
			zap.String("workflow_name", cb.Workflow.WorkflowName))
	}
	return nil
}

func (z *zadigImpl) domainMonitor(cb model.Callback) error {
	switch cb.Workflow.Status {
	case "passed":
		return z.handleDomainMonitorPassed()
	default:
		defaultStage := cb.Workflow.Stages[0]
		domains := make([]string, 0)
		for _, j := range defaultStage.Jobs {
			displayName := j.DisplayName
			status := j.Status
			if status == "failed" {
				domains = append(domains, displayName)
			}
		}
		return z.handleDomainMonitorDefault(domains)
	}
}

func (z *zadigImpl) deploySubEnv(cb model.Callback) error {
	switch cb.Workflow.Status {
	case "passed":
		return z.handleDeploySubEnvPassed(cb)
	default:
		defaultStage := cb.Workflow.Stages[0]
		domains := make([]string, 0)
		for _, j := range defaultStage.Jobs {
			displayName := j.DisplayName
			status := j.Status
			if status == "failed" {
				domains = append(domains, displayName)
			}
		}
		return z.handleDomainMonitorDefault(domains)
	}
}

func (z *zadigImpl) handleDomainMonitorPassed() error {
	req := model.SendInteractiveReq{
		TemplateId:  "ctp_AAzXWvvEaFd5",
		Target:      model.Group,
		ReceiveName: "Devops Notification",
		Params: map[string]string{
			"title":   "域名监控运行成功",
			"content": "所有域名运行正常，且没有在 30 天内到期的域名，无需任何操作。",
		},
	}
	return z.lark.SendInteractive(req)
}

func (z *zadigImpl) handleDeploySubEnvPassed(cb model.Callback) error {
	totalTime := cb.Workflow.EndTime - cb.Workflow.CreateTime
	detail, err := z.GetTaskDetail("test33", cb.Workflow.TaskID)
	if err != nil {
		return err
	}
	var subEnv string
	for _, param := range detail.Params {
		if param.Name == "环境" {
			subEnv = param.Value
		}
	}
	var services, branches []string
	for _, stage := range detail.Stages {
		if stage.Name == "构建" {
			for _, job := range stage.Jobs {
				services = append(services, job.JobInfo.ServiceName)
				for _, repo := range job.Spec.Repos {
					branches = append(branches, repo.Branch)
				}
			}
		}
	}
	req := model.SendInteractiveReq{
		TemplateId:  "ctp_AAz7KWuUUkkh",
		Target:      model.User,
		ReceiveName: cb.Workflow.TaskCreatorEmail,
		Params: map[string]string{
			"project_name":    cb.Workflow.ProjectName,
			"workflow_name":   "fat-base-workflow",
			"workflow_number": strconv.Itoa(cb.Workflow.TaskID),
			"duration":        fmt.Sprintf("%02d:%02d", totalTime/60, totalTime%60),
			"host":            "",
			"sub_env":         subEnv,
			"service":         strings.Join(services, "\n"),
			"branch":          strings.Join(branches, "\n"),
		},
	}
	return z.lark.SendInteractive(req)
}

func (z *zadigImpl) handleDomainMonitorDefault(domains []string) error {
	content := ""
	for _, domain := range domains {
		content += "请关注：" + domain + "\n"
	}
	req := model.SendInteractiveReq{
		TemplateId:  "ctp_AAzXaSRdmtsX",
		Target:      model.Group,
		ReceiveName: "Engineering Incident Report Group",
		Params: map[string]string{
			"title":   "域名即将过期提醒",
			"content": content,
		},
	}
	return z.lark.SendInteractive(req)
}
