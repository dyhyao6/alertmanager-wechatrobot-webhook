package model

// FeiShuMessage represents a FeiShu (Lark) message
type FeiShuMessage struct {
	MsgType string            `json:"msg_type"`
	Content *FeiShuContent    `json:"content"`
}

type FeiShuContent struct {
	Text string `json:"text,omitempty"`
}

// FeiShuMarkdown represents FeiShu markdown content
type FeiShuMarkdown struct {
	Text string `json:"text"`
}