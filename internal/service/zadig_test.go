package service

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jeremy2566/octopipe/internal/model"
	"go.uber.org/zap"
	"resty.dev/v3"
)

func TestZadigImpl_GetTestEnvList(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)
	client.SetDebug(true)
	z := NewZadig(log, client)
	envs, _ := z.GetTestEnvList("fat-base-envrionment")
	for _, env := range envs {
		log.Info(fmt.Sprintf("%+v", env))
	}
}

func TestZadigImpl_GetTestEnvDetail(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)

	z := NewZadig(log, client)
	_, _ = z.GetTestEnvDetail("test34", "fat-base-envrionment")
}

func TestZadigImpl_selectedNamespace(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)

	z := NewZadig(log, client).(*zadigImpl)
	for i := 0; i < 100; i++ {
		println(z.selectedNamespace())
	}
}

func TestCreate_Namespace_E2E(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)

	z := NewZadig(log, client).(*zadigImpl)
	subEnv, _ := z.CreateSubEnv()
	taskId, _ := z.DeployService(model.DeployServiceReq{
		SubEnv:      subEnv,
		ServiceName: "backoffice-v1-web",
		BranchName:  "feature/INF-666",
		GithubActor: "jeremy2566",
	})
	log.Info("deploy service.", zap.Int("taskId", taskId))
	_ = z.AddService(model.AddServiceReq{
		SubEnv:      subEnv,
		ServiceName: "payment-api",
	})
	taskId, _ = z.DeployService(model.DeployServiceReq{
		SubEnv:      subEnv,
		ServiceName: "payment-api",
		BranchName:  "feature/INF-666",
		GithubActor: "jeremy2566",
	})
}

func TestZadigImpl_handleDomainMonitorDefault(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)

	z := NewZadig(log, client).(*zadigImpl)
	_ = z.handleDomainMonitorDefault([]string{"storehub.com"})
}

func TestZadigImpl_handleDomainMonitorPassed(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)

	z := NewZadig(log, client).(*zadigImpl)
	_ = z.handleDomainMonitorPassed()
}

func TestZadigImpl_GetServiceCharts(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)

	z := NewZadig(log, client).(*zadigImpl)
	charts := z.GetServiceCharts()
	for k, v := range charts {
		log.Info(fmt.Sprintf("%s: %s", k, v))
	}
}

func TestZadigImpl_Webhook01(t *testing.T) {
	callback := `{
    "object_kind": "workflow",
    "event": "workflow",
    "workflow": {
        "task_id": 61,
        "project_name": "devops-tools",
        "project_display_name": "DevOps_Tools",
        "workflow_name": "domain-monitor",
        "workflow_display_name": "domain-monitor",
        "status": "failed",
        "remark": "",
        "detail_url": "http://zadigx.shub.us/v1/projects/detail/devops-tools/pipelines/custom/domain-monitor?display_name=domain-monitor",
        "error": "",
        "create_time": 1757387902,
        "start_time": 1757387904,
        "end_time": 1757387968,
        "stages": [
            {
                "name": "default",
                "status": "failed",
                "start_time": 1757387904,
                "end_time": 1757387968,
                "jobs": [
                    {
                        "name": "job-0-0-0-shub-us",
                        "display_name": "shub-us",
                        "type": "freestyle",
                        "status": "failed",
                        "start_time": 1757387904,
                        "end_time": 1757387909,
                        "error": "waitJobStart: pod failed, jobName:domain-monitor-61-krkw2, podName:domain-monitor-61-krkw2-7d4xt\nconditions info: type:PodReadyToStartContainers, status:False, reason:, message:\ntype:Initialized, status:True, reason:, message:\ntype:Ready, status:False, reason:PodFailed, message:\ntype:ContainersReady, status:False, reason:PodFailed, message:\ntype:PodScheduled, status:True, reason:, message:\n",
                        "spec": {
                            "repositories": null,
                            "image": ""
                        }
                    },
                    {
                        "name": "job-0-1-0-storehubhq-com",
                        "display_name": "storehubhq-com",
                        "type": "freestyle",
                        "status": "passed",
                        "start_time": 1757387904,
                        "end_time": 1757387925,
                        "error": "",
                        "spec": {
                            "repositories": null,
                            "image": ""
                        }
                    },
                    {
                        "name": "job-0-2-0-bpit-me",
                        "display_name": "bpit-me",
                        "type": "freestyle",
                        "status": "failed",
                        "start_time": 1757387904,
                        "end_time": 1757387908,
                        "error": "waitJobStart: pod failed, jobName:domain-monitor-61-lvskg, podName:domain-monitor-61-lvskg-9x89c\nconditions info: type:PodReadyToStartContainers, status:False, reason:, message:\ntype:Initialized, status:True, reason:, message:\ntype:Ready, status:False, reason:PodFailed, message:\ntype:ContainersReady, status:False, reason:PodFailed, message:\ntype:PodScheduled, status:True, reason:, message:\n",
                        "spec": {
                            "repositories": null,
                            "image": ""
                        }
                    },
                    {
                        "name": "job-0-3-0-mymyhub-com",
                        "display_name": "mymyhub-com",
                        "type": "freestyle",
                        "status": "passed",
                        "start_time": 1757387904,
                        "end_time": 1757387968,
                        "error": "",
                        "spec": {
                            "repositories": null,
                            "image": ""
                        }
                    },
                    {
                        "name": "job-0-4-0-beepit-com",
                        "display_name": "beepit-com",
                        "type": "freestyle",
                        "status": "passed",
                        "start_time": 1757387904,
                        "end_time": 1757387960,
                        "error": "",
                        "spec": {
                            "repositories": null,
                            "image": ""
                        }
                    },
                    {
                        "name": "job-0-5-0-storehub-com",
                        "display_name": "storehub-com",
                        "type": "freestyle",
                        "status": "passed",
                        "start_time": 1757387904,
                        "end_time": 1757387954,
                        "error": "",
                        "spec": {
                            "repositories": null,
                            "image": ""
                        }
                    }
                ],
                "error": ""
            }
        ],
        "task_creator": "hanzhang",
        "task_creator_id": "97882c5f-a266-11ef-aa9f-02058eeea235",
        "task_creator_phone": "13325666101",
        "task_creator_email": "jeremy.zhang@storehub.com",
        "task_type": "workflow"
    }
}`
	var cb model.Callback
	json.Unmarshal([]byte(callback), &cb)
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)

	z := NewZadig(log, client).(*zadigImpl)

	_ = z.Webhook(cb)

}
