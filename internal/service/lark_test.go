package service

//
//import (
//	"fmt"
//	"testing"
//	"time"
//
//	"github.com/jeremy2566/octopipe/internal/model"
//	"go.uber.org/zap"
//	"resty.dev/v3"
//)
//
//func TestHandler_Sender_Failed(t *testing.T) {
//	log, _ := zap.NewDevelopment()
//	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)
//	client.SetDebug(true)
//	l := NewLark(log, client)
//	totalTime := 19
//
//	req := model.SenderReq{
//		ProjectName:    "fat-base-environment",
//		WorkflowName:   "fat-base-workflow",
//		WorkflowNumber: 5357,
//		Duration:       fmt.Sprintf("%02d:%02d", totalTime/60, totalTime%60),
//		SubEnv:         "test34",
//		Service:        "payment-api",
//		Branch:         "feature/INF-666",
//		Success:        false,
//		Email:          "jeremy.zhang@storehub.com",
//	}
//	l.Sender(req)
//}
//
//func TestHandler_Sender_Success(t *testing.T) {
//	log, _ := zap.NewDevelopment()
//	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)
//	l := NewLark(log, client)
//
//	totalTime := 19
//
//	req := model.SenderReq{
//		ProjectName:    "fat-base-environment",
//		WorkflowName:   "fat-base-workflow",
//		WorkflowNumber: 5357,
//		Duration:       fmt.Sprintf("%02d:%02d", totalTime/60, totalTime%60),
//		SubEnv:         "test34",
//		Service:        "payment-api",
//		Branch:         "feature/INF-666",
//		Success:        true,
//		Email:          "jeremy.zhang@storehub.com",
//	}
//	l.Sender(req)
//}
//
//func TestHandler_Sender_Failed_Namespace(t *testing.T) {
//	log, _ := zap.NewDevelopment()
//	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)
//	l := NewLark(log, client)
//
//	req := model.SenderReq{
//		ProjectName:    "fat-base-environment",
//		WorkflowName:   "fat-base-workflow",
//		WorkflowNumber: 0,
//		Duration:       "",
//		SubEnv:         "test17",
//		Service:        "",
//		Branch:         "feature/INF-666",
//		Success:        false,
//		Email:          "jeremy.zhang@storehub.com",
//	}
//	l.Sender(req)
//}
//
//func TestHandler_DomainMonitor(t *testing.T) {
//	req := model.SenderLarkReq{
//		ReceiveIdType: "chat_id",
//		ReceiveId:     "oc_b97835f507ec3d6648a3445bdf85d549",
//		MsgType:       "interactive",
//		Content: model.ContentLarkReq{
//			Type: "template",
//			Data: model.DataLarkReq{
//				TemplateID: "ctp_AAzXWvvEaFd5",
//			},
//		},
//	}
//	log, _ := zap.NewDevelopment()
//	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)
//	l := NewLark(log, client)
//	l.DomainMonitor(req)
//}

import (
	"fmt"
	"testing"
	"time"

	"github.com/jeremy2566/octopipe/internal/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"resty.dev/v3"
)

func TestLarkImpl_getDevOpsBotToken(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)
	l := NewLark(log, client)

	larkImpl := l.(*larkImpl)
	token, err := larkImpl.getDevOpsBotToken()
	assert.Nil(t, err)
	log.Info(token)
}

func TestLarkImpl_getAllGroupForDevOpsBot(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)
	l := NewLark(log, client)

	larkImpl := l.(*larkImpl)
	groups := larkImpl.getAllGroupForDevOpsBot()
	for k, v := range groups {
		fmt.Println(k, "=>", v)
	}
}

func TestLarkImpl_SendInteractive_To_Group_Success(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)
	client.SetDebug(true)
	l := NewLark(log, client)

	req := model.SendInteractiveReq{
		TemplateId:  "ctp_AAzXWvvEaFd5",
		Target:      model.Group,
		ReceiveName: "Lark Bot",
		Params: map[string]string{
			"title":   "域名监控运行成功",
			"content": "所有域名运行正常，且没有在 30 天内到期的域名，无需任何操作。",
		},
	}
	err := l.SendInteractive(req)
	assert.Nil(t, err)
}

func TestLarkImpl_SendInteractive_To_User_Success(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)
	client.SetDebug(true)
	l := NewLark(log, client)

	req := model.SendInteractiveReq{
		TemplateId:  "ctp_AAz7KWuUUkkh",
		Target:      model.User,
		ReceiveName: "jeremy.zhang@storehub.com",
		Params: map[string]string{
			"project_name": "域名监控运行成功",
			"content":      "所有域名运行正常，且没有在 30 天内到期的域名，无需任何操作。",
		},
	}
	err := l.SendInteractive(req)
	assert.Nil(t, err)
}

func TestLarkImpl_SendInteractive_To_User_Fail(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)
	client.SetDebug(true)
	l := NewLark(log, client)

	req := model.SendInteractiveReq{
		TemplateId:  "ctp_AAzXaSRdmtsX",
		Target:      model.User,
		ReceiveName: "jeremy.zhang@storehub.com",
		Params: map[string]string{
			"title":   "域名过期提醒",
			"content": "请检查 storehub.com",
		},
	}
	err := l.SendInteractive(req)
	assert.Nil(t, err)
}

func TestLarkImpl_SendInteractive_To_Group_Fail(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := resty.New().SetRetryCount(3).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second)
	client.SetDebug(true)
	l := NewLark(log, client)

	req := model.SendInteractiveReq{
		TemplateId:  "ctp_AAzXaSRdmtsX",
		Target:      model.Group,
		ReceiveName: "Lark Bot",
		Params: map[string]string{
			"title":   "域名过期提醒",
			"content": "请检查 storehub.com",
		},
	}
	err := l.SendInteractive(req)
	assert.Nil(t, err)
}
