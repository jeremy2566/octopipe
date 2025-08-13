package service

import (
	"fmt"
	"testing"
	"time"

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
