package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (h Handler) GetTenantToken() (string, error) {
	body := struct {
		AppID     string `json:"app_id"`
		AppSecret string `json:"app_secret"`
	}{
		AppID:     "cli_a6ea84eab9f9902f",
		AppSecret: "YCIgBYAdtPwQfYpt7oALlfn0fw2IdPnK",
	}

	resp := struct {
		Code              int    `json:"code"`
		TenantAccessToken string `json:"tenant_access_token"`
	}{}

	res, err := h.client.R().
		SetBody(body).
		SetContentType("application/json").
		SetResult(&resp).
		Post("https://open.larksuite.com/open-apis/auth/v3/tenant_access_token/internal")
	if err != nil {
		return "", fmt.Errorf("get tenant token: %w", err)
	}
	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("get tenant token status: %s", res.String())
	}

	return resp.TenantAccessToken, nil
}

type TemplateVariable struct {
	Duration       string `json:"duration"`
	Host           string `json:"host"`
	ProjectName    string `json:"project_name"`
	WorkflowName   string `json:"workflow_name"`
	WorkflowNumber string `json:"workflow_number"`
	SubEnv         string `json:"sub_env"`
	Service        string `json:"service"`
	Branch         string `json:"branch"`
}

type Data struct {
	TemplateID       string           `json:"template_id"`
	TemplateVariable TemplateVariable `json:"template_variable"`
}

type Content struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}

type Message struct {
	ReceiveID string `json:"receive_id"`
	MsgType   string `json:"msg_type"`
	Content   string `json:"content"`
}

func (h Handler) Sender(req SenderReq) {
	tv := TemplateVariable{
		Duration:       req.Duration,
		Host:           req.Host,
		ProjectName:    "fat-base-environment",
		WorkflowName:   "fat-base-workflow",
		WorkflowNumber: req.WorkflowNumber,
		SubEnv:         req.SubEnv,
		Service:        req.Service,
		Branch:         req.Branch,
	}
	var tempId string
	if req.Success {
		tempId = "ctp_AAz7KWuUUkkh"
	} else {
		tempId = "ctp_AAz7KDoTWE2h"
	}
	data := Data{
		TemplateID:       tempId,
		TemplateVariable: tv,
	}
	c := Content{
		Type: "template",
		Data: data,
	}

	cJson, err := json.Marshal(c)
	if err != nil {
		h.log.Warn("send lark msg failed.", zap.Error(err))
		return
	}
	body := Message{
		ReceiveID: req.Email,
		MsgType:   "interactive",
		Content:   string(cJson),
	}

	marshal, _ := json.Marshal(body)
	fmt.Println(string(marshal))
	resp := struct {
		msg string `json:"msg"`
	}{}
	token, _ := h.GetTenantToken()
	res, err := h.client.R().
		SetBody(body).
		SetResult(&resp).
		SetAuthToken(token).
		Post("https://open.larksuite.com/open-apis/im/v1/messages?receive_id_type=email")
	if err != nil {

	}

	fmt.Println(res)

}

type SenderReq struct {
	Duration       string
	Host           string
	WorkflowNumber string
	SubEnv         string
	Service        string
	Branch         string
	Success        bool
	Email          string
}
