package model

type SenderReq struct {
	ProjectName    string
	WorkflowName   string
	WorkflowNumber int
	Duration       string
	SubEnv         string
	Service        string
	Branch         string
	Success        bool
	Email          string
}

type DataReq struct {
	TemplateID       string              `json:"template_id"`
	TemplateVariable TemplateVariableReq `json:"template_variable"`
}

type ContentReq struct {
	Type string  `json:"type"`
	Data DataReq `json:"data"`
}

type MessageReq struct {
	ReceiveID string `json:"receive_id"`
	MsgType   string `json:"msg_type"`
	Content   string `json:"content"`
}

type TemplateVariableReq struct {
	Duration       string `json:"duration"`
	Host           string `json:"host"`
	ProjectName    string `json:"project_name"`
	WorkflowName   string `json:"workflow_name"`
	WorkflowNumber int    `json:"workflow_number"`
	SubEnv         string `json:"sub_env"`
	Service        string `json:"service"`
	Branch         string `json:"branch"`
}

type DomainMonitorReq struct {
	Success   bool
	TmplId    string
	ReceiveId string // 群聊 id
	params    map[string]string
}

type SenderLarkReq struct {
	ReceiveIdType string
	ReceiveId     string
	MsgType       string
	Content       ContentLarkReq
}

type DataLarkReq struct {
	TemplateID string `json:"template_id"`
}

type ContentLarkReq struct {
	Type string      `json:"type"`
	Data DataLarkReq `json:"data"`
}
