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
		Host:           "https://zadigx.shub.us/v1/projects/detail/fat-base-envrionment/pipelines/custom/test33/5354?status=failed&id=&display_name=fat-base-workflow",
		WorkflowNumber: "5354",
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

	req := SenderReq{
		Duration:       "2m30",
		Host:           "https://zadigx.shub.us/v1/projects/detail/fat-base-envrionment/pipelines/custom/test33/5357?status=passed&id=&display_name=fat-base-workflow",
		WorkflowNumber: "5357",
		SubEnv:         "test34",
		Service:        "payment-api",
		Branch:         "feature/INF-666",
		Success:        true,
		Email:          "jeremy.zhang@storehub.com",
	}
	h.Sender(req)
}
