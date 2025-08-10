package handler

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
	"resty.dev/v3"
)

func TestHandler_GetTenantToken(t *testing.T) {
	h := Handler{
		log:    zap.NewNop(),
		client: resty.New(),
	}

	token, _ := h.GetTenantToken()
	fmt.Println(token)
}

func TestHandler_Sender_Failed(t *testing.T) {
	h := Handler{
		log:    zap.NewNop(),
		client: resty.New(),
	}

	req := SenderReq{
		Duration:       "2m31",
		WorkflowNumber: 5354,
		SubEnv:         "test34",
		Service:        "payment-api",
		Branch:         "feature/INF-666",
		Success:        false,
		Email:          "jeremy.zhang@storehub.com",
	}
	h.Sender(req)
}

func TestHandler_Sender_Success(t *testing.T) {
	h := Handler{
		log:    zap.NewNop(),
		client: resty.New(),
	}

	totalTime := 19

	req := SenderReq{
		ProjectName:    "fat-base-environment",
		WorkflowName:   "fat-base-workflow",
		WorkflowNumber: 5357,
		Duration:       fmt.Sprintf("%02d:%02d", totalTime/60, totalTime%60),
		SubEnv:         "test34",
		Service:        "payment-api",
		Branch:         "feature/INF-666",
		Success:        true,
		Email:          "jeremy.zhang@storehub.com",
	}
	h.Sender(req)
}
