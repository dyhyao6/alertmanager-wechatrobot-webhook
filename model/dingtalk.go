package model

// DingTalkMessage represents a DingTalk message
type DingTalkMessage struct {
	MsgType string       `json:"msgtype"`
	Markdown *DingTalkMarkdown `json:"markdown,omitempty"`
	Text    *DingTalkText `json:"text,omitempty"`
}

type DingTalkMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingTalkText struct {
	Content string `json:"content"`
}