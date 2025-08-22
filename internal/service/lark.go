package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeremy2566/octopipe/internal/model"
	"go.uber.org/zap"
	"resty.dev/v3"
)

var _ Lark = &larkImpl{}

type Lark interface {
	Sender(req model.SenderReq)
	DomainMonitor(req model.SenderLarkReq) error
}

type larkImpl struct {
	log    *zap.Logger
	client *resty.Client
}

type sender struct {
	ReceiveId string `json:"receiveId"`
	MsgType   string `json:"msg_type"`
	Content   string `json:"content"`
}

func (l *larkImpl) DomainMonitor(req model.SenderLarkReq) error {
	token, err := l.GetTenantToken()
	if err != nil {
		return err
	}
	content, _ := json.Marshal(req.Content)

	body := sender{
		ReceiveId: req.ReceiveId,
		MsgType:   req.MsgType,
		Content:   string(content),
	}

	resp, err := l.client.R().
		SetBody(body).
		SetAuthToken(token).
		Post(fmt.Sprintf("/open-apis/im/v1/messages?receive_id_type=%s", req.ReceiveIdType))
	fmt.Println(resp)
	return nil
}

func (l *larkImpl) Sender(req model.SenderReq) {
	var tempId, status string
	if req.Success {
		tempId = "ctp_AAz7KWuUUkkh"
		status = "passed"
	} else {
		tempId = "ctp_AAz7KDoTWE2h"
		status = "failed"
	}
	host := fmt.Sprintf(
		"https://zadigx.shub.us/v1/projects/detail/%s/pipelines/custom/test33/%d?status=%s&id=&display_name=%s",
		req.ProjectName,
		req.WorkflowNumber,
		status,
		req.WorkflowName,
	)
	tv := model.TemplateVariableReq{
		Duration:       req.Duration,
		Host:           host,
		ProjectName:    req.ProjectName,
		WorkflowName:   req.WorkflowName,
		WorkflowNumber: req.WorkflowNumber,
		SubEnv:         req.SubEnv,
		Service:        req.Service,
		Branch:         req.Branch,
	}

	data := model.DataReq{
		TemplateID:       tempId,
		TemplateVariable: tv,
	}
	c := model.ContentReq{
		Type: "template",
		Data: data,
	}

	cJson, err := json.Marshal(c)
	if err != nil {
		l.log.Warn("send lark msg failed.", zap.Error(err))
		return
	}
	body := model.MessageReq{
		ReceiveID: req.Email,
		MsgType:   "interactive",
		Content:   string(cJson),
	}

	marshal, _ := json.Marshal(body)
	fmt.Println(string(marshal))
	resp := struct {
		msg string `json:"msg"`
	}{}
	token, _ := l.GetTenantToken()
	res, err := l.client.R().
		SetBody(body).
		SetResult(&resp).
		SetAuthToken(token).
		Post("/open-apis/im/v1/messages?receive_id_type=email")
	if err != nil {

	}

	fmt.Println(res)

}

func NewLark(log *zap.Logger, client *resty.Client) Lark {
	client.SetBaseURL("https://open.larksuite.com")
	return &larkImpl{
		log:    log,
		client: client,
	}
}

func (l *larkImpl) GetTenantToken() (string, error) {
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

	res, err := l.client.R().
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
