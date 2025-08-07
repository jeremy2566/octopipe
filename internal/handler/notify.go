package handler

type Body struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (h Handler) Notify(msg string) {
	body := Body{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: msg,
		},
	}

	_, err := h.client.R().
		SetBody(body).
		Post("https://open.larksuite.com/open-apis/bot/v2/hook/57d28364-4710-49e1-a8e4-eea29c82ab49")
	if err != nil {
		h.Notify("notify failed")
	}
}
