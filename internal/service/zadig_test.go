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
        "task_id": 5732,
        "project_name": "fat-base-envrionment",
        "project_display_name": "fat-base-environment",
        "workflow_name": "test33",
        "workflow_display_name": "fat-base-workflow",
        "status": "passed",
        "remark": "",
        "detail_url": "http://zadigx.shub.us/v1/projects/detail/fat-base-envrionment/pipelines/custom/test33?display_name=fat-base-workflow",
        "error": "",
        "create_time": 1757401639,
        "start_time": 1757401640,
        "end_time": 1757401834,
        "stages": [
            {
                "name": "构建",
                "status": "passed",
                "start_time": 1757401640,
                "end_time": 1757401834,
                "jobs": [
                    {
                        "name": "job-1-0-0-构建发布",
                        "display_name": "构建发布-backoffice-v1-web-app-backoffice-v1-web",
                        "type": "zadig-build",
                        "status": "passed",
                        "start_time": 1757401640,
                        "end_time": 1757401834,
                        "error": "",
                        "spec": {
                            "repositories": [
                                {
                                    "source": "github",
                                    "repo_owner": "storehubnet",
                                    "repo_namespace": "storehubnet",
                                    "repo_name": "backoffice-v1-web",
                                    "branch": "release",
                                    "prs": null,
                                    "tag": "",
                                    "author_name": "",
                                    "commit_id": "0bee75101106e3eecc7825c34fc34e1a41195104",
                                    "commit_url": "https://github.com/storehubnet/backoffice-v1-web/commit/0bee7510",
                                    "commit_message": "Merge pull request #4240 from storehubnet/CB-14652-from-release\n\nfeat(CB-14652): Add GrowthBook whitelist control for promotions v2 page"
                                }
                            ],
                            "image": "858157298152.dkr.ecr.ap-southeast-1.amazonaws.com/backoffice-v1-web:release-0bee7510"
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

func TestZadigImpl_Webhook02(t *testing.T) {
	callback := `{
    "object_kind": "workflow",
    "event": "workflow",
    "workflow": {
        "task_id": 5734,
        "project_name": "fat-base-envrionment",
        "project_display_name": "fat-base-environment",
        "workflow_name": "test33",
        "workflow_display_name": "fat-base-workflow",
        "status": "passed",
        "remark": "",
        "detail_url": "http://zadigx.shub.us/v1/projects/detail/fat-base-envrionment/pipelines/custom/test33?display_name=fat-base-workflow",
        "error": "",
        "create_time": 1757403005,
        "start_time": 1757403006,
        "end_time": 1757403163,
        "stages": [
            {
                "name": "构建",
                "status": "failed",
                "start_time": 1757403006,
                "end_time": 1757403163,
                "jobs": [
                    {
                        "name": "job-1-0-0-构建发布",
                        "display_name": "构建发布-backoffice-v1-web-app-backoffice-v1-web",
                        "type": "zadig-build",
                        "status": "failed",
                        "start_time": 1757403006,
                        "end_time": 1757403023,
                        "error": "",
                        "spec": {
                            "repositories": [
                                {
                                    "source": "github",
                                    "repo_owner": "storehubnet",
                                    "repo_namespace": "storehubnet",
                                    "repo_name": "backoffice-v1-web",
                                    "branch": "k8s-pro-test",
                                    "prs": null,
                                    "tag": "",
                                    "author_name": "",
                                    "commit_id": "eb2bf5924c1c2b8b6916c9b4ff57405e12db449c",
                                    "commit_url": "https://github.com/storehubnet/backoffice-v1-web/commit/eb2bf592",
                                    "commit_message": "Edit domain"
                                }
                            ],
                            "image": ""
                        }
                    },
                    {
                        "name": "job-1-0-1-构建发布",
                        "display_name": "构建发布-core-api-core-api",
                        "type": "zadig-build",
                        "status": "failed",
                        "start_time": 1757403006,
                        "end_time": 1757403163,
                        "error": "",
                        "spec": {
                            "repositories": [
                                {
                                    "source": "github",
                                    "repo_owner": "storehubnet",
                                    "repo_namespace": "storehubnet",
                                    "repo_name": "core-api",
                                    "branch": "develop",
                                    "prs": null,
                                    "tag": "",
                                    "author_name": "",
                                    "commit_id": "50bdf7f4b40ba126f6ab2d49dd3e482e49bf1079",
                                    "commit_url": "https://github.com/storehubnet/core-api/commit/50bdf7f4",
                                    "commit_message": "Merge pull request #587 from storehubnet/autoupdatesubmodule\n\nauto update submodule"
                                }
                            ],
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

func TestZadigImpl_GetTaskDetail(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)

	z := NewZadig(log, client)
	resp, _ := z.GetTaskDetail("test33", 5731)
	fmt.Println(resp)
}

func TestZadigImpl_handleDeploySubEnvPassed(t *testing.T) {
	callback := `{
    "object_kind": "workflow",
    "event": "workflow",
    "workflow": {
        "task_id": 5734,
        "project_name": "fat-base-envrionment",
        "project_display_name": "fat-base-environment",
        "workflow_name": "test33",
        "workflow_display_name": "fat-base-workflow",
        "status": "passed",
        "remark": "",
        "detail_url": "http://zadigx.shub.us/v1/projects/detail/fat-base-envrionment/pipelines/custom/test33?display_name=fat-base-workflow",
        "error": "",
        "create_time": 1757403005,
        "start_time": 1757403006,
        "end_time": 1757403163,
        "stages": [
            {
                "name": "构建",
                "status": "failed",
                "start_time": 1757403006,
                "end_time": 1757403163,
                "jobs": [
                    {
                        "name": "job-1-0-0-构建发布",
                        "display_name": "构建发布-backoffice-v1-web-app-backoffice-v1-web",
                        "type": "zadig-build",
                        "status": "failed",
                        "start_time": 1757403006,
                        "end_time": 1757403023,
                        "error": "",
                        "spec": {
                            "repositories": [
                                {
                                    "source": "github",
                                    "repo_owner": "storehubnet",
                                    "repo_namespace": "storehubnet",
                                    "repo_name": "backoffice-v1-web",
                                    "branch": "k8s-pro-test",
                                    "prs": null,
                                    "tag": "",
                                    "author_name": "",
                                    "commit_id": "eb2bf5924c1c2b8b6916c9b4ff57405e12db449c",
                                    "commit_url": "https://github.com/storehubnet/backoffice-v1-web/commit/eb2bf592",
                                    "commit_message": "Edit domain"
                                }
                            ],
                            "image": ""
                        }
                    },
                    {
                        "name": "job-1-0-1-构建发布",
                        "display_name": "构建发布-core-api-core-api",
                        "type": "zadig-build",
                        "status": "failed",
                        "start_time": 1757403006,
                        "end_time": 1757403163,
                        "error": "",
                        "spec": {
                            "repositories": [
                                {
                                    "source": "github",
                                    "repo_owner": "storehubnet",
                                    "repo_namespace": "storehubnet",
                                    "repo_name": "core-api",
                                    "branch": "develop",
                                    "prs": null,
                                    "tag": "",
                                    "author_name": "",
                                    "commit_id": "50bdf7f4b40ba126f6ab2d49dd3e482e49bf1079",
                                    "commit_url": "https://github.com/storehubnet/core-api/commit/50bdf7f4",
                                    "commit_message": "Merge pull request #587 from storehubnet/autoupdatesubmodule\n\nauto update submodule"
                                }
                            ],
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

	_ = z.handleDeploySubEnvPassed(cb)
}
