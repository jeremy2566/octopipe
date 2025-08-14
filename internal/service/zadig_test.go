package service

import (
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
	z.AddService(model.AddServiceReq{
		SubEnv:      subEnv,
		ServiceName: "payment-api",
	})
	z.DeployService(model.DeployServiceReq{
		SubEnv:      subEnv,
		ServiceName: "payment-api",
		BranchName:  "feature/INF-666",
		GithubActor: "jeremy2566",
	})
}
